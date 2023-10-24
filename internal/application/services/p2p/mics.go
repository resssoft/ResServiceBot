package p2p

import "strconv"

func Number(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}
