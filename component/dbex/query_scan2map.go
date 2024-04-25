package dbex

import (
	"context"
	"database/sql"
)

// QueryScan2Map 执行SQL查询并将结果转换为map[string]interface{}格式的切片。
// ctx: 上下文，用于控制查询的取消、超时等。
// db: 用于执行查询的SQL数据库实例。
// sqlstr: 待执行的SQL查询语句。
// args: SQL查询语句中使用的参数。
// 返回值: 查询结果转换后的map切片和可能出现的错误。
func QueryScan2Map(ctx context.Context, db *sql.DB, sqlstr string, args ...interface{}) ([]map[string]interface{}, error) {
	// 准备查询语句
	stmt, err := db.PrepareContext(ctx, sqlstr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close() // 确保查询语句被关闭

	// 执行查询
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 确保查询结果集被关闭

	// 获取查询返回的列名
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 计算列数，并初始化结果集和存放列值的切片
	var (
		colCount = len(cols)
		result   = []map[string]interface{}{}
		values   = make([]interface{}, colCount)
		valPtrs  = make([]interface{}, colCount)
	)

	// 遍历所有行，并将每行数据转换为map格式
	for rows.Next() {
		for i := 0; i < colCount; i++ {
			valPtrs[i] = &values[i]
		}
		rows.Scan(valPtrs...)

		entry := map[string]interface{}{}
		for i, col := range cols {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b) // 将字节切片转换为字符串
			} else {
				v = val
			}
			entry[col] = v
		}

		result = append(result, entry)
	}

	return result, nil
}
