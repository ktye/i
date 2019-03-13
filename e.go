package i

// Eval
func (a *A) E(l interface{}) interface{} {
	return evl(a, l)
}

func evl(v ...interface{}) interface{} {
	_ = v[0].(*A)            // a←
	_ = v[1].([]interface{}) // l←
	return td()
}
