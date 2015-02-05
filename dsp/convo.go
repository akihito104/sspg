package dsp

func Convolve(x []int16, y []int) []int {
	res := make([]int, len(x)+len(y)-1)
	for i, s := range x {
		ss := int(s)
		for j, p := range y {
			res[i+j] += ss * p
		}
	}
	return res
}
