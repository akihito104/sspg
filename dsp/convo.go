package dsp

func Convolve(x []int16, y []int32) []int32 {
	res := make([]int32, len(x)+len(y)-1)
	for i, s := range x {
		ss := int32(s)
		for j, p := range y {
			res[i+j] += ss * p
		}
	}
	return res
}
