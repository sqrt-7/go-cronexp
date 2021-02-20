package cronexp

// fill fills an []int slice within the given range
func fill(min, max int) []int {
	res := make([]int, max-min+1)
	for i := 0; i <= max-min; i++ {
		res[i] = i + min
	}

	return res
}
