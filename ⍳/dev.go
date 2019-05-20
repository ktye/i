package main

// devvars(device variables)
// ⍳ defines: iv(k-version), io(read,write)
// others: im(image:pixel buffer), pl(ot), ed(itor) planned for github.com/ktye/ui/examples/interpret

func iv(x, y v) v { // k-version. ⍳ returns 1
	if x == nil {
		return 1
	}
	return nil
}

func devvars(a map[v]v) { // register devvars in k-tree
	a["devvars"] = map[s]func(v, v) v{
		"iv": iv,
	}
}
