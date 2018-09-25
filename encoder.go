package main

import "github.com/gushitong/aryadb/io"

const (
	MaxKeySize = 255
)

type Symbol byte

const (
	SymbolSet    Symbol = 's'
	SymbolHash   Symbol = 'h'
	SymbolList   Symbol = 'l'
	SymbolZset   Symbol = 'z'
	SymbolString Symbol = 'k'
	SymbolMinSeq Symbol = 'a'
	SymbolMaxSeq Symbol = 'b'
)

type aryCommand struct {
	Raw  []byte
	Args [][]byte
}

func (aryCmd aryCommand) StringKey() []byte {
	return EStringKey(aryCmd.Args[0])
}

func (aryCmd aryCommand) HashKey() ([]byte, error) {
	return EHashKey(aryCmd.Args[0], aryCmd.Args[1])
}

// string key encode: {SymbolString} + key
func EStringKey(key []byte) []byte {
	return io.ConcatBytearray(
		[]byte{(byte)(SymbolString)},
		key,
	)
}

// hash key encode: {SymbolHash} + {len(key)} + key + hashKey
func EHashKey(key, hashKey []byte)([]byte, error) {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}

	return io.ConcatBytearray(
		[]byte{(byte)(SymbolHash)},
		[]byte{uint8(len(key))},
		key,
		hashKey,
	), nil
}

// set key encode: {SymbolSet} + {len(key)} + key + setKey
func ESetKey(key, setKey []byte)([]byte, error) {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}

	return io.ConcatBytearray(
		[]byte{(byte)(SymbolSet)},
		[]byte{uint8(len(key))},
		key,
		setKey,
	), nil
}

// zset key encode: {SymbolZset} + {len(key)} + key + len(zsetKey) + zsetKey
func EZsetKey(key, zsetKey []byte)([]byte, error) {
	if len(key) > MaxKeySize || len(zsetKey) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}
	return io.ConcatBytearray(
		[]byte{(byte)(SymbolZset)},
		[]byte{uint8(len(key))},
		key,
		[]byte{uint8(len(zsetKey))},
		zsetKey,
	), nil
}

// zset score encode: {SymbolZset} + len(key) + {symbol} + zsetScore
func EZsetScore(key, score []byte)([]byte, error) {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}
	s, err := io.ParseInt64(score)
	if err != nil{
		return nil, ErrIntegerValue
	}
	symbol := byte('=')
	if s < 0 {
		symbol = byte('-')
	}

	return io.ConcatBytearray(
		[]byte{(byte)(SymbolZset)},
		[]byte{uint8(len(key))},
		key,
		[]byte{symbol},
		score,
	), nil
}

// list key encode: {SymbolList} + {len(key)} + key + seq
func EListKey(key, seq []byte) ([]byte, error)  {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}

	return io.ConcatBytearray(
		[]byte{(byte)(SymbolList)},
		[]byte{uint8(len(key))},
		key,
		seq,
	), nil
}

// list key min-seq encode: {SymbolList} + {len(key)} + key + {SymbolMinSeq}
func EListMinSeq(key []byte) ([]byte, error) {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}

	return io.ConcatBytearray(
		[]byte{(byte)(SymbolList)},
		[]byte{uint8(len(key))},
		key,
		[]byte{(byte)(SymbolMinSeq)},
	), nil
}

// list key max-seq encode: {SymbolList} + {len(key)} + key + {SymbolMaxSeq}
func EListMaxSeq(key []byte) ([]byte, error) {
	if len(key) > MaxKeySize {
		return nil, ErrKeySizeExceeded
	}

	return io.ConcatBytearray(
		[]byte{(byte)(SymbolList)},
		[]byte{uint8(len(key))},
		key,
		[]byte{(byte)(SymbolMaxSeq)},
	), nil
}

// encode hash prefix
func EHashPrefix(key []byte)([]byte, error) {
	if len(key) > 255 {
		return nil, ErrKeySizeExceeded
	}

	return io.ConcatBytearray(
		[]byte{(byte)(SymbolHash)},
		[]byte{uint8(len(key))},
		key,
	), nil
}

// decode hash key
func DHashKey(hashKey []byte)([]byte, []byte, error) {
	if hashKey[0] != (byte)(SymbolHash) {
		return nil, nil, ErrCorruptedHashKey
	}
	lenKey := int(hashKey[1])
	if len(hashKey) <= lenKey+2 {
		return nil, nil, ErrCorruptedHashKey
	}
	return hashKey[2:lenKey+2], hashKey[lenKey+2:], nil
}