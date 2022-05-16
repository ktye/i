package swtch

func f(x int32) int32 {
	switch x {
	case 0:
		return (1 + x)
	case 1:
		return (2 * x)
	default:
		return x
	}
}
