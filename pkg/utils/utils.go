package utils

import "strconv"

const (
	_parserBase = 10
	_parserSize = 64
)

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, _parserBase)
}
