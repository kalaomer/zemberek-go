package data

// WeightLookup is an interface for weight lookup operations
type WeightLookup interface {
	Get(key string) float32
	Size() int
}
