package loop

func f() {
	var i int32
	i = int32(0)
	for ; i < 3; i = (i + 1) {
		i = (i * 2)
		continue
	}
}
