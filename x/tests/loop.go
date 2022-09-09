package loop

func f() {
	var n, i int32
	i = 0
	for ; i < 3; i = (i + 1) {
		n = (n + 1)
	}
}
