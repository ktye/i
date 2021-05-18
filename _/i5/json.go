package k

import (
	"encoding/json"
	"fmt"
)

//? `json?"1 " /1
func decodeJson(x T) T {
	c, o := x.(C)
	if o == false {
		panic("type")
	}
	dx(x)

	var f interface{}
	e := json.Unmarshal(c.v, &f)
	if e != nil {
		panic(e)
	}
	return kj(f)
}

//@ `json 1 2 3 /"[1,2,3]"
//@ `json(1;2 3;4.) /"[1,[2,3],4]"
func encodeJson(x T) T {
	defer dx(x)
	b, e := json.Marshal(x)
	if e != nil {
		panic(e)
	}
	return KC(b)
}

func kj(f interface{}) T {
	switch v := f.(type) {
	case bool:
		return v
	case float64:
		if float64(int(v)) == v {
			return int(v)
		}
		return v
	case string:
		return KC([]byte(v))
	case nil:
		return nil
	case []interface{}:
		r := make([]T, len(v))
		for i := range r {
			r[i] = kj(v[i])
		}
		return KL(r).uf()
	case map[string]interface{}:
		key := make([]string, len(v))
		val := make([]T, len(v))
		i := 0
		for s, u := range v {
			key[i] = s
			val[i] = kj(u)
			i++
		}
		return dict(KS(key), KL(val).uf())
	default:
		fmt.Printf("kj %T %v\n", v, v)
		panic("type")
	}
}
