package main

import "time"

func regtime() timestamp {
	return timestamp{Time: time.Now()}
}

type timestamp struct {
	time.Time
}

func (t timestamp) cpy() timestamp {
	return t
}

func (t timestamp) String() string {
	return t.Format("2006.01.02T15:04:05")
}

// ConvertTo is used to convert values to time by dyadic $
// t$0        / current time
// t$20340293 / seconds since y2k
// t$"2019.04.04T21:22:33" / from string
// t$2019 04 04 / from a numeric vector ([]complex128)
// By default "t" stores a time value.
func (t timestamp) ConvertTo(u v) v {
	if s, o := u.(s); o {
		tm, err := time.Parse(s, "2006.01.02T15:04:05")
		if err != nil {
			panic(err)
		}
		return timestamp{tm}
	} else if z, o := u.(complex128); o {
		sec := real(z)
		if sec == 0 {
			return timestamp{time.Now()}
		}
		return timestamp{y2k.Add(time.Duration(float64(time.Second) * sec))}
	} else if vec, o := u.([]complex128); o {
		w := [8]int{0, 1, 1, 0, 0, 0, 0, 0}
		for i := range vec {
			w[i] = int(real(vec[i]))
		}
		return timestamp{time.Date(w[0], time.Month(w[1]), w[2], w[3], w[4], w[5], w[6], nil)}
	} else {
		panic("type")
	}
}

var y2k time.Time

func init() {
	y2k, _ = time.Parse("2006.01.02", "2000.01.01")
}
