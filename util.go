package main

import (
	"math"
)

type Util struct {
}

func (util *Util) MinInt(a int, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
