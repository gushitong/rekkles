package main

import (
	"strings"
)

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

func RegisterHandler(s *server, key string, handler Handler) {
	s.Handlers[key] = handler
}
