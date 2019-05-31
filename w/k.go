package w

type b = byte
type k = uint32
type i = int32
type f = float64

//                i  f   z  s  c  g  l  d
var lns = [9]k{0, 4, 8, 16, 4, 1, 4, 4, 8}
var e k = 0xFFFFFFFF

var m []b

func mk(t b, n int) k { // make type t of len n (-1:atom)
	sz := lns[t]
	if n >= 0 {
		sz *= k(n)
	}
	sz += 8 // size needed including header
	bs := k(16)
	bt := 0
	for i := 4; i < 30; i++ { // calculate bucket bt from size sz (clz)
		if sz <= bs {
			bt = i
			break
		}
		bs <<= 1
	}
	if bt == 0 {
		return e
	}
	fb, a := 0, k(0)
	for i := bt; i < 30; i++ { // find next free bucket >= bt
		if k(i) >= get(4) {
			grw()
		}
		if get(k(4*i)) != 0 {
			fb, a = i, get(k(4*i))
			break
		}
	}
	if fb == 0 {
		return e
	}
	for i := fb - 1; i >= bt; i-- { // split large buckets
		put(k(4*i), a)
		m[a] = b(i)
		a += k(1) << b(i)
		m[a] = b(i)
	}
	if n < 0 { // set header
		m[int(a+1)] = t
	} else {
		put(k(a), k(m[int(a)]|t<<5)|k(n)<<8)
	}
	put(a+4, 1) // refcount
	return a
}
func typ(a k) (b, int) { // type and length at addr
	i := int(a)
	t := m[i] >> 5
	if t == 0 {
		return m[int(i+1)], -1
	}
	return t, int(get(k(i)) >> 8) //int(m[i+1]) | int(m[i+2]<<8) | int(m[i+3]<<16)
}
func rst() { // reset memory
	m = make([]b, 1<<16)
	p := k(len(m))
	for i := 15; i > 6; i-- {
		p >>= 1
		m[p] = b(i)
		put(k(4*i), p)
	}
	m[0] = 7
	put(4, 16) // total memory (log2)
	// TODO: pointer to k-tree at 8
	put(k(4*9), 0)   // no free bucket 9
	put(1<<9, k(73)) // 73: 1<<6|9 (type i, bucket 9), length is ignored
	for i := range lns {
		put(k(4*i+8)+1<<9, k(lns[i]))
	}
}
func put(a, x k) {
	i := int(a)
	m[i] = b(x)
	m[i+1] = b(x >> 8)
	m[i+2] = b(x >> 16)
	m[i+3] = b(x >> 24)
}
func get(a k) k { i := int(a); return k(m[i]) | k(m[i+1])<<8 | k(m[i+2])<<16 | k(m[i+3])<<24 }
func grw() {
	s := m[4]
	if 1<<k(s) != len(m) {
		panic("grw")
	}
	put(k(4*s), k(len(m)))
	m = append(m, make([]b, len(m))...)
	m[4] = s + 1
	m[1<<s] = s
}
