package main

// fluid simulation based on nigel's implementation in golang.org/x/exp/shiny/example/fluid
// jos stam, "Real-Time Fluid Dynamics for Games"
//
// in the interpreter:
//  .m:{(*x 1)flow(*x 2)}  /register mouse callback
//  \L10:0 flow 0          /start the animation
//                         /click on the image

const (
	NN         = 128
	iterations = 20
	dt         = 0.1
	diff       = 0
	visc       = 0
	force      = 5
	source     = 20
	fade       = 0.89
)

type array [NN + 2][NN + 2]float32

var (
	dens, densPrev array
	u, uPrev       array
	v, vPrev       array
	xPrev, yPrev   int
	havePrevLoc    bool
	pix            [4 * NN * NN]byte
)

func flw(x, y k) (r k) {
	xx, yy := m.k[2+x], m.k[2+y]
	if xx > NN {
		xx = 0
	}
	if yy > NN {
		yy = 0
	}
	step(int(xx), int(yy))
	return drw(mki(NN), inc(flow))
}
func step(x, y int) { // x, y mouse clicks
	if x+y > 0 {
		dens[x][y] = source
		if havePrevLoc {
			u[x+1][y+1] = force * float32(x-xPrev)
			v[x+1][y+1] = force * float32(y-yPrev)
		}
		xPrev, yPrev, havePrevLoc = x, y, true
	}
	velStep(&u, &v, &uPrev, &vPrev)
	densStep(&dens, &densPrev, &u, &v)
	for i := range dens {
		for j := range dens[i] {
			dens[i][j] *= fade
		}
	}
	pp := 8 + flow<<2
	for y := 0; y < NN; y++ {
		for x := 0; x < NN; x++ {
			d := int32(dens[x+1][y+1] * 0xff)
			if d < 0 {
				d = 0
			} else if d > 0xff {
				d = 0xff
			}
			v := 255 - uint8(d)
			p := k(NN*y+x) * 4
			m.c[pp+p+0] = v
			m.c[pp+p+1] = v
			m.c[pp+p+2] = v
			m.c[pp+p+3] = 0xff
		}
	}
}
func addSource(x, s *array) {
	for i := range x {
		for j := range x[i] {
			x[i][j] += dt * s[i][j]
		}
	}
}
func setBnd(b int, x *array) {
	switch b {
	case 0:
		for i := 1; i <= NN; i++ {
			x[0+0][i] = x[1][i]
			x[NN+1][i] = x[NN][i]
			x[i][0+0] = x[i][1]
			x[i][NN+1] = x[i][NN]
		}
	case 1:
		for i := 1; i <= NN; i++ {
			x[0+0][i] = -x[1][i]
			x[NN+1][i] = -x[NN][i]
			x[i][0+0] = x[i][1]
			x[i][NN+1] = x[i][NN]
		}
	case 2:
		for i := 1; i <= NN; i++ {
			x[0+0][i] = x[1][i]
			x[NN+1][i] = x[NN][i]
			x[i][0+0] = -x[i][1]
			x[i][NN+1] = -x[i][NN]
		}
	}
	x[0+0][0+0] = 0.5 * (x[1][0+0] + x[0+0][1])
	x[0+0][NN+1] = 0.5 * (x[1][NN+1] + x[0+0][NN])
	x[NN+1][0+0] = 0.5 * (x[NN][0+0] + x[NN+1][1])
	x[NN+1][NN+1] = 0.5 * (x[NN][NN+1] + x[NN+1][NN])
}
func linSolve(b int, x, x0 *array, a, c float32) {
	if a == 0 && c == 1 {
		for i := 1; i <= NN; i++ {
			for j := 1; j <= NN; j++ {
				x[i][j] = x0[i][j]
			}
		}
		setBnd(b, x)
		return
	}
	invC := 1 / c
	for k := 0; k < iterations; k++ {
		for i := 1; i <= NN; i++ {
			for j := 1; j <= NN; j++ {
				x[i][j] = (x0[i][j] + a*(x[i-1][j]+x[i+1][j]+x[i][j-1]+x[i][j+1])) * invC
			}
		}
		setBnd(b, x)
	}
}
func diffuse(b int, x, x0 *array, diff float32) {
	a := dt * diff * NN * NN
	linSolve(b, x, x0, a, 1+4*a)
}
func advect(b int, d, d0, u, v *array) {
	const dt0 = dt * NN
	for i := 1; i <= NN; i++ {
		for j := 1; j <= NN; j++ {
			x := float32(i) - dt0*u[i][j]
			if x < 0.5 {
				x = 0.5
			}
			if x > NN+0.5 {
				x = NN + 0.5
			}
			i0 := int(x)
			i1 := i0 + 1

			y := float32(j) - dt0*v[i][j]
			if y < 0.5 {
				y = 0.5
			}
			if y > NN+0.5 {
				y = NN + 0.5
			}
			j0 := int(y)
			j1 := j0 + 1

			s1 := x - float32(i0)
			s0 := 1 - s1
			t1 := y - float32(j0)
			t0 := 1 - t1
			d[i][j] = s0*(t0*d0[i0][j0]+t1*d0[i0][j1]) + s1*(t0*d0[i1][j0]+t1*d0[i1][j1])
		}
	}
	setBnd(b, d)
}
func project(u, v, p, div *array) {
	for i := 1; i <= NN; i++ {
		for j := 1; j <= NN; j++ {
			div[i][j] = (u[i+1][j] - u[i-1][j] + v[i][j+1] - v[i][j-1]) / (-2 * NN)
			p[i][j] = 0
		}
	}
	setBnd(0, div)
	setBnd(0, p)
	linSolve(0, p, div, 1, 4)
	for i := 1; i <= NN; i++ {
		for j := 1; j <= NN; j++ {
			u[i][j] -= (NN / 2) * (p[i+1][j+0] - p[i-1][j+0])
			v[i][j] -= (NN / 2) * (p[i+0][j+1] - p[i+0][j-1])
		}
	}
	setBnd(1, u)
	setBnd(2, v)
}
func velStep(u, v, u0, v0 *array) {
	addSource(u, u0)
	addSource(v, v0)
	u0, u = u, u0
	diffuse(1, u, u0, visc)
	v0, v = v, v0
	diffuse(2, v, v0, visc)
	project(u, v, u0, v0)
	u0, u = u, u0
	v0, v = v, v0
	advect(1, u, u0, u0, v0)
	advect(2, v, v0, u0, v0)
	project(u, v, u0, v0)
}
func densStep(x, x0, u, v *array) {
	addSource(x, x0)
	x0, x = x, x0
	diffuse(0, x, x0, diff)
	x0, x = x, x0
	advect(0, x, x0, u, v)
}
