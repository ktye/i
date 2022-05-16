package label

func f() {
	var i int32
L:
	for {
		i += 1
		if i > 3 {
			break L
		}
	}
}
func g() {
	var i int32
	for {
		i += 1
		if i > 3 {
			break
		}
	}
}
