package group

type ISection interface {
	GetBeg() int64
	GetEnd() int64
}

func Section(beg, end, step int64) [][]int64 {
	sections := [][]int64{}
	for i := beg; i < end; i += step {
		e := i + step
		if e > end {
			e = end
		}
		sections = append(sections, []int64{i, e})
	}
	return sections
}
func SectionAny[T ISection](root ISection, step int64, f func(int64, int64) T) []T {
	beg := root.GetBeg()
	end := root.GetEnd()
	sections := []T{}
	for i := beg; i < end; i += step {
		e := i + step
		if e > end {
			e = end
		}
		s := f(i, e)
		sections = append(sections, s)
	}
	return sections
}
