package main

import "github.com/gushitong/rekkles/ut"

type HashEncoder struct {
	key []byte
}

func (h HashEncoder) Encode(field []byte) []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolHash)},
		[]byte{uint8(len(h.key))},
		h.key,
		field,
	)
}

func (h HashEncoder) Prefix() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolHash)},
		[]byte{uint8(len(h.key))},
		h.key,
	)
}

func (h HashEncoder) QueueKey() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolQueue)},
		h.key,
	)
}

func (h HashEncoder) Decode(hashKey []byte) ([]byte, error) {
	if hashKey[0] != (byte)(SymbolHash) {
		return nil, ErrCorruptedHashKey
	}
	lenKey := int(hashKey[1])
	if len(hashKey) <= lenKey+2 {
		return nil, ErrCorruptedHashKey
	}
	return hashKey[lenKey+2:], nil
}

func NewHashEncoder(key []byte) (*HashEncoder, error) {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}
	return &HashEncoder{key: key}, nil
}
