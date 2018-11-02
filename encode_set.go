package main

import "github.com/gushitong/rekkles/ut"

type SetEncoder struct {
	key []byte
}

func (s SetEncoder) EncodeMember(member []byte) []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolSet)},
		[]byte{uint8(len(s.key))},
		s.key,
		member,
	)
}

func (s SetEncoder) QueueKey() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolQueue)},
		s.key,
	)
}

func (s SetEncoder) Prefix() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolSet)},
		[]byte{uint8(len(s.key))},
		s.key,
	)
}

func (s SetEncoder) DecodeMember(setKey []byte) ([]byte, error) {
	if len(setKey) < 3 {
		return nil, ErrCorruptedSetKey
	}
	lKey := int(setKey[1])
	return setKey[lKey+2:], nil
}

func NewSetEncoder(key []byte) (*SetEncoder, error) {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}
	return &SetEncoder{key: key}, nil
}
