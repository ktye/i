package k

func first(x K) K {
	n := cnt(x)
	if n < 0 {
		return x
	}
	return ati(x, 0)
}
