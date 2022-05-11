package cast

func f(x int32) int64   { return int64(x) }
func g(x int64) int32   { return int32(x) }
func h(x int32) float64 { return float64(x) }
func i(x float64) int32 { return int32(x) }
func j(x int32) uint32  { return uint32(x) }
func k(x uint32) int32  { return int32(x) }
