package idsnow

func offsetSum(ls ...int) int64 {
	t := int64(0)
	for _, l := range ls {
		t += int64(l)
	}
	return t
}
