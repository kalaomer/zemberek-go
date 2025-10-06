package hash

// Mphf is an interface for Minimum Perfect Hash Functions
type Mphf interface {
	Get(key interface{}, initialHash ...int32) int32
}

// Rshift performs an unsigned right shift operation
func Rshift(val int32, n uint) int32 {
	return int32(uint32(val) >> n)
}
