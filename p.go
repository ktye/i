package i

// Parse
func (a *A) P(s string) interface{} {
	return prs(a, s)
}

func prs(v ...interface{}) interface{} {
	_ = v[0].(*A)     // a←
	_ = v[1].(string) // s←
	return td()
}
