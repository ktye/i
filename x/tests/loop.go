package loop

func f() {
	var n, i int32
	n = 0
	i = 0
	for ; i < 3; i = (i + 1) {
		n = (n + 1)
	}
}
