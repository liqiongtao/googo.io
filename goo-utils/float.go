package goo_utils

type Float64 float64

func (f Float64) ToFixed(n int) float64 {
	if n == 0 {
		return float64(int64(f))
	}

	base := int64(10)
	for i := 1; i < n; i++ {
		base = base * 10
	}

	v := int64(float64(f) * float64(base))
	return float64(v) / float64(base)
}

func (f Float64) ToPercent(n int) float64 {
	if n == 0 {
		n = 2
	}

	base := int64(10)
	for i := 1; i < n; i++ {
		base = base * 10
	}

	v := int64(float64(f*100) * float64(base))
	return float64(v) / float64(base)
}
