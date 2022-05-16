package cont

func f() {
	var i int32
	for {
		i += 1
		if i < 2 {
			continue
		}
		i *= 2
	}
}
