package io

import (
	"strings"
	"strconv"
	"fmt"
	"encoding/json"
)

func DumpMap(v map[string]string) {
	js, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(js))
}

func ClearBit(b byte, pos uint) byte {
	return b &^ (1 << pos)
}

func SetBit(b byte, pos uint) byte {
	return b |^ (1 << pos)
}

func SliceIndex(l, o int) int {
	if l == 0 {
		return 0
	}
	if o < -1*l {
		o = 0
	}else if o < 0 {
		o += l
	}else if o >= l {
		o = l - 1
	}
	return o
}

func LowerString(s []byte) string {
	return strings.ToLower(string(s))
}

func ParseInt64(v []byte) (int64, error) {
	return strconv.ParseInt(string(v), 10, 64)
}

func ParseFloat64(v []byte) (float64, error) {
	return strconv.ParseFloat(string(v), 64)
}

func Int642Byte(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}

func Float642Byte(v float64) []byte {
	return []byte(strconv.FormatFloat(v, 'f', -1, 64))
}

func ConcatBytearray(bytesarray ...[]byte) []byte {
	var totalLen int
	for _, s := range bytesarray {
		totalLen += len(s)
	}
	bytes := make([]byte, totalLen)

	var i int
	for _, s := range bytesarray {
		i += copy(bytes[i:], s)
	}
	return bytes
}