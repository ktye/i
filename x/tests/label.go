package label

func f() {
	var i int32
	i = 0
L:
	for {
		i = (i + 1)
		if i > 3 {
			break L
		}
	}
}
func g() {
	var i int32
	i = 0
	for {
		i = (i + 1)
		if i > 3 {
			break
		}
	}
}
