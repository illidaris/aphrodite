package idsnow

func Offset(lens []int, index int) int64 {
	l := len(lens)
	if l == 0 {
		return 0
	}
	if index > l-1 {
		return 0
	}
	return offsetSum(lens[index:]...)
}

func IdPartsFrmVals(lens []int, vals ...int64) []int64 {
	parts := []int64{}
	for i := 0; i < len(lens); i++ {
		val := int64(0)
		if i < len(vals) {
			val = vals[i]
		}
		offset := Offset(lens, i+1)
		parts = append(parts, val<<offset)
	}
	return parts
}
func GetValsFrmId(lens []int, id int64) []int64 {
	vals := []int64{}
	for i := 0; i < len(lens); i++ {
		offset := Offset(lens, i+1)
		maskSequence := (int64(1)<<int64(lens[i]) - int64(1)) << offset
		vals = append(vals, (id&maskSequence)>>offset)
	}
	return vals
}

func offsetSum(ls ...int) int64 {
	t := int64(0)
	for _, l := range ls {
		t += int64(l)
	}
	return t
}
