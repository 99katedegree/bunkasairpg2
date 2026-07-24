package battle

// RNG is a mulberry32 pseudo-random number generator.
// Algorithm matches the TypeScript implementation in the frontend exactly.
type RNG struct {
	s uint32
}

func NewRNG(seed int64) *RNG {
	return &RNG{s: uint32(seed)}
}

func (r *RNG) Next() float64 {
	r.s += 0x6d2b79f5
	t := uint32(int32(r.s^(r.s>>15)) * int32(1|r.s))
	t = (t + uint32(int32(t^(t>>7))*int32(61|t))) ^ t
	return float64(t^(t>>14)) / 4294967296.0
}
