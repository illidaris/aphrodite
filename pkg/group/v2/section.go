package v2

type ISection interface {
	GetBeg() int64
	GetEnd() int64
}

/*
Section 将区间 [beg, end) 按指定步长分割为连续子区间
参数:
  - beg : 起始位置（包含）
  - end : 结束位置（不包含）
  - step: 每个子区间的步长值

返回值:
  - 二维切片，每个元素表示一个子区间 [start, end]，最后一个子区间会自动调整结束边界
*/
func Section(beg, end, step int64) [][]int64 {
	sections := [][]int64{}
	// 循环生成区间切片，步长控制分割跨度
	for i := beg; i < end; i += step {
		e := i + step
		// 处理最后一个子区间的边界溢出
		if e > end {
			e = end
		}
		sections = append(sections, []int64{i, e})
	}
	return sections
}

/*
SectionAny 将原始区间分割为泛型实例集合
类型参数:
  - T : 必须实现 ISection 接口的类型

参数:
  - root: 原始区间对象，需实现 ISection 接口
  - step: 每个子区间的步长值
  - f   : 接收子区间起止位置并生成对应泛型实例的回调函数

返回值:
  - 泛型实例切片，每个元素对应分割后的子区间
*/
func SectionAny[T ISection](root ISection, step int64, f func(int64, int64) T) []T {
	beg := root.GetBeg()
	end := root.GetEnd()
	sections := []T{}
	// 通过泛型工厂函数生成区间实例
	for i := beg; i < end; i += step {
		e := i + step
		// 边界检查保证区间完整性
		if e > end {
			e = end
		}
		s := f(i, e)
		sections = append(sections, s)
	}
	return sections
}
