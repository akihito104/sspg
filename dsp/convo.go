package dsp

func Convolve(x []int16, y []int32) []int64 {
	res := make([]int64, len(x)+len(y)-1)
	for i, s := range x {
		ss := int64(s)
		for j, p := range y {
			res[i+j] += ss * int64(p)
		}
	}
	return res
}
