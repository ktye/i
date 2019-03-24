package i

type rn = rune
type rv = []rn
type sf func(rv) i

func prs(v interface{}) interface{} { // string contains no comments
	b, n := rv(v.(s)), 0
	for _, r := range b { // trim left
		if !any(r, wsp) {
			break
		}
		n++
	}
	b = b[n:]
	var c rv
	for { // preserve strings, disambiguate minus, replace \n
		if len(b) == 0 {
			break
		} else if len(b) == 1 {
			c = append(c, b[0])
			break
		}
		if n := sStr(b); n > 0 {
			c = append(c, b[:n]...)
			b = b[n:]
		} else if n := sNum(b); n > 0 {
			c = append(c, b[:n]...)
			b = b[n:]
		} else if b[0] == '\r' {
			b = b[1:]
		} else if b[0] == '\n' {
			c = append(c, ';')
			b = b[1:]
		} else if b[0] == '-' {
			c = append(c, '-', ' ')
			b = b[1:]
		}
	}
	if len(c) == 0 {
		return l{}
	}
	for i := len(c) - 1; i >= 0; i-- { // trim right
		if !any(c[i], wsp) {
			break
		}
		c = c[:i]
	}
	if len(c) == 0 {
		return l{}
	}
	var r l
	c, r = pLst(c, nil, false)
	if len(c) != 0 {
		return e("parse")
	}
	return r
}

func pLst(b rv, term sf, cull bool) (rv, l) {
	/*
		r := l{}
		for {
			if len(b) < 1 {
				break
			}
			if term != nil {
				if n := term(b); n > 0 {
					b = b[:n]
					break
				}
			}
			for {
				if len(b) > 0
			}
		}
	*/
	e("nyi")
	return nil, nil
}

const dig = "0123456789"
const con = "πø"
const sym = `+\-*%!&|<>=~,^#_$?@.`
const uni = `⍉×÷⍳⍸⌊⌽⌈⍋⌸≡∧⍴≢↑⌊↓⍕∪⍎⍣¯ℜℑ√⍟`
const uav = "⍨¨⌿⍀"
const wsp = " \t\r"

func sNum(s rv) i {
	sn := 0
	for i, r := range s {
		switch {
		case i == 0 && !any(r, "-+"+dig):
			return 0
		case any(r, dig):
		case any(r, "eEaj"):
			if i == 1 && any(s[0], "-+") {
				return 0
			}
			if n := sNum(s[i:]); n > 0 {
				return i + n
			}
			return i
		default:
			break
		}
		sn++
	}
	if sn == 1 && any(s[0], "+-") {
		return 0
	}
	return sn
}
func sNam(s rv) i {
	a := func(r rn) bool {
		if alpha(r) {
			return true
		}
		return false
	}
	n := 0
	for i, r := range s {
		switch {
		case i == 0 && any(r, con):
			return 1
		case i == 0 && !a(r):
			return 0
		case a(r) || any(r, "0123456789"):
		default:
			return i
		}
		n++
	}
	return n
}
func sSym(s rv) i {
	if s[0] != '`' {
		return 0
	}
	if len(s) == 1 {
		return 1
	}
	return 1 + sNam(s[1:])
}
func sStr(s rv) i {
	if len(s) < 2 || s[0] != '"' {
		return 0
	}
	h := false
	for i, r := range s {
		switch {
		case i == 0:
		case r == '\\':
			h = !h
		case r == '"' && !h:
			return i
		}
	}
	return 0
}
func sVrb(s rv) i {
	for _, r := range s {
		if any(r, sym) || any(r, uni) {
			return 1
		}
		return 0
	}
	return 0
}
func sAsn(s rv) i {
	if n := sVrb(s); n != 0 && len(s) > n && s[n] == ':' {
		return n + 1
	}
	return 0
}
func sIov(s rv) i {
	if len(s) < 2 {
		return 0
	}
	if any(s[0], dig) && s[1] == ':' {
		return 2
	}
	return 0
}
func sAdv(s rv) i {
	for i, r := range s {
		if i == 0 && any(r, uav) {
			return 1
		}
	}
	if any(s[0], `'/\`) {
		if len(s) > 1 && s[1] == ':' {
			return 2
		}
		return 1
	}
	return 0
}
func sSem(s rv) i { return pref(s, ";") }
func sCol(s rv) i { return pref(s, ":") }
func sViw(s rv) i { return pref(s, "::") }
func sCnd(s rv) i { return pref(s, `$[`) }
func sDct(s rv) i {
	if pref(s, "[") == 0 {
		return 0
	}
	for i, r := range s[1:] {
		switch {
		case alpha(r):
		case i > 0 && i == ':':
			return i
		default:
			break
		}
	}
	return 0
}
func sObr(s rv) i { return pref(s, "[") }
func sOpa(s rv) i { return pref(s, "(") }
func sOcb(s rv) i { return pref(s, "{") }
func sCbr(s rv) i { return pref(s, "]") }
func sCpa(s rv) i { return pref(s, ")") }
func sCcb(s rv) i { return pref(s, "}") }

func any(r rn, s s) bool {
	for _, x := range s {
		if r == x {
			return true
		}
	}
	return false
}
func pref(r rv, p string) int {
	s := string(r)
	if len(s) < len(p) {
		return 0
	}
	if s[:len(p)] == p {
		return 1
	}
	return 0
}
func alpha(r rn) bool {
	if (r >= 'a' && r <= 'a') || (r >= 'A' && r <= 'Z') {
		return true
	}
	return false
}
