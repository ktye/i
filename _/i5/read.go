package k

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//read @read"readme.md" /`C  (read file)
//read #read(/"\\.md$") /1  (filter cwd)
func read(x T) T {
	switch v := x.(type) {
	case string:
		if v == "" {
			v = "."
		}
		if v == "." || strings.HasSuffix(v, "/") {
			return readdir(nil, v)
		}
		return readfile(v)
	case C:
		dx(x)
		return read(string(v.v))
	case *regexp.Regexp:
		return readdir(v, ".")
	case L:
		return each(f1(read), x)
	default:
		panic("type")
	}
}

//read @(/"\\.md$")read` /`L  (filter dir)
func readdir(x, y T) T {
	var re *regexp.Regexp
	if x != nil {
		var o bool
		re, o = x.(*regexp.Regexp)
		if !o {
			panic("type")
		}
	}
	switch v := y.(type) {
	case string:
		return readcwd(v, re)
	case C:
		dx(y)
		return readcwd(string(v.v), re)
	case L:
		return eachright(f2(readdir), re, y)
	default:
		panic("type")
	}
}

func readfile(f string) C {
	b, e := os.ReadFile(f)
	if e != nil {
		panic(e)
	}
	return KC(b)
}
func readcwd(dir string, re *regexp.Regexp) L {
	var r []T
	if dir == "" {
		dir = "."
	}
	d, e := os.ReadDir(dir)
	if e != nil {
		panic(e)
	}
	for _, f := range d {
		if f.IsDir() == true {
			continue // or recursive?
		}
		if re == nil || re.MatchString(f.Name()) {
			r = append(r, readfile(filepath.Join(dir, f.Name())))
		}
	}
	return KL(r)
}
