package main

import (
	"math"
	"strconv"
)

type Util struct {
}

func (util *Util) MinInt(a int, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func (util *Util) GetInt(flag *Flag, defaultValue int) int {
	if flag == nil {
		return defaultValue
	}

	ret, err := strconv.Atoi(flag.Value)
	if err != nil {
		return defaultValue
	}

	return ret
}

func (util *Util) GetBoolean(flag *Flag, defaultValue bool) bool {
	if flag == nil {
		return defaultValue
	}

	ret, err := strconv.ParseBool(flag.Value)
	if err != nil {
		return defaultValue
	}

	return ret
}
