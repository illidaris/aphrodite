package elasticex

import (
	"time"
)

type BatchSpan struct {
	IsSkip  bool
	Size    int64
	Anchors []interface{} // 锚点
}

func NeedPagingSpan(need, batch int64, isSkip bool) []*BatchSpan {
	ps := []*BatchSpan{}
	// 防止死循环
	timeout := time.After(time.Minute * 30)
	for need > 0 {
		isBreak := false
		select {
		case <-timeout:
			isBreak = true
		default:
			if need > batch {
				ps = append(ps, &BatchSpan{IsSkip: isSkip, Size: int64(batch)})
			} else {
				ps = append(ps, &BatchSpan{IsSkip: isSkip, Size: int64(need)})
			}
			need -= batch
		}
		if isBreak {
			break
		}
	}
	return ps
}

type LimitSearch struct {
	Limit int64 // 窗口搜索量限制

	Offset int64        // 偏移量，从多少点位开始获取数据
	Total  int64        // 总计量，总计获取数据量
	Spans  []*BatchSpan // 搜索分段
}

func (l *LimitSearch) InitSpans() {
	ps := []*BatchSpan{}
	ps = append(ps, NeedPagingSpan(l.Offset, l.Limit, true)...) // 需要跳过的
	ps = append(ps, NeedPagingSpan(l.Total, l.Limit, false)...) // 需要获取的
	l.Spans = ps
}

// RepairSpans 根据实际数据优化分段，减少无效查询次数
func (l *LimitSearch) RepairSpans(realTotal int64) {
	need := realTotal
	if need <= 0 {
		l.Spans = []*BatchSpan{}
	}
	cut := need / l.Limit
	if need%l.Limit != 0 {
		cut++
	}
	if cut < int64(len(l.Spans)) {
		l.Spans = l.Spans[:cut]
	}
}

func (l *LimitSearch) SearchByStep(f func(index, subCursor, subSize int64) ([]interface{}, int64, error)) ([]interface{}, int64, error) {
	l.InitSpans()
	// 先搜索一次
	if len(l.Spans) == 0 {
		return nil, 0, nil
	}
	result := []interface{}{}
	firstRows, total, err := f(0, 0, l.Spans[0].Size)
	if err != nil {
		return result, total, err
	}
	if !l.Spans[0].IsSkip {
		result = append(result, firstRows...)
	}
	l.RepairSpans(total)
	for index, span := range l.Spans {
		if index == 0 {
			continue
		}
		rows, _, err := f(int64(index), 0, l.Spans[0].Size)
		if err != nil {
			return result, total, err
		}

		if !span.IsSkip {
			result = append(result, rows...)
		}
	}
	return result, total, err
}
