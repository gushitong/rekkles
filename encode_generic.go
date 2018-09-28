package main

import "github.com/gushitong/aryadb/ut"

type Symbol byte

const (
	SymbolSet    Symbol = 's'
	SymbolHash   Symbol = 'h'
	SymbolZset   Symbol = 'z'
	SymbolString Symbol = 'k'

	SymbolList      Symbol = 'l'
	SymbolDelimiter Symbol = '='
	SymbolMinSeq    Symbol = 'a'
	SymbolMaxSeq    Symbol = 'b'
	SymbolQueue Symbol = 'q' // 容量
)

const (
	MaxKeySize = 255
)

var SetMemberValue = []byte{uint8(1)}

type StringEncoder struct {
	key []byte
}

func (s StringEncoder) Encode() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolString)},
		s.key,
	)
}

func NewStringEncoder(key []byte) *StringEncoder {
	return &StringEncoder{key:key}
}