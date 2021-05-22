package k

func Neg(x K) K { return Ki(-int32(x)) }

func Sqr(x K) K { trap(Nyi); return x }

func Add(x, y K) K { return Ki(int32(x) + int32(y)) }
func Sub(x, y K) K { return Ki(int32(x) - int32(y)) }
func Mul(x, y K) K { return Ki(int32(x) * int32(y)) }
func Div(x, y K) K { return Ki(int32(x) / int32(y)) }

func Min(x, y K) K { trap(Nyi); return x }
func Max(x, y K) K { trap(Nyi); return x }
