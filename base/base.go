package base

import (
	"math"
	"strings"
)

func Uint64ToString(i uint64, baseStr string) string {
	var s []rune
	baseLen := uint64(len(baseStr))
	for i > 0 {
		round := i / baseLen
		remain := i % baseLen
		s = append(s, rune(baseStr[remain]))
		i = round
	}
	return string(s)
}

func StringToUint64(s string, baseStr string) uint64 {
	var i uint64
	baseLen := len(baseStr)
	for idx, r := range s {
		base := uint64(math.Pow(float64(baseLen), float64(idx)))
		i += uint64(strings.Index(baseStr, string(r))) * base
	}
	return i
}
