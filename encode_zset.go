package main

import "github.com/gushitong/aryadb/ut"

type ZsetEncoder struct {
	key []byte
}

func (z ZsetEncoder) EncodeMemberKey(member []byte) []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolSet)},
		[]byte{uint8(len(z.key))},
		z.key,
		[]byte{uint8(len(member))},
		member,
	)
}

func (z ZsetEncoder) MemberPrefix() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolSet)},
		[]byte{uint8(len(z.key))},
		z.key,
	)
}

func (z ZsetEncoder) EncodeScoreKey(score int64) []byte {
	var flag byte
	if score < 0 {
		flag = '-'
	} else {
		flag = '='
	}
	zscore := ut.Int642Bytes(score)
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolZset)},
		[]byte{uint8(len(z.key))},
		z.key,
		[]byte{flag},
		zscore,
	)
}

func (z ZsetEncoder) ScorePrefix() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolZset)},
		[]byte{uint8(len(z.key))},
		z.key,
	)
}

func (z ZsetEncoder) EncodeScore(score []byte) ([]byte, error) {
	s, err := ut.ParseInt64(score)
	if err != nil {
		return nil, err
	}
	return ut.Int642Bytes(s), nil
}

func (z ZsetEncoder) DecodeScoreKey(scoreKey []byte) (int64, error) {
	if len(scoreKey) < 5 || scoreKey[0] != (byte)(SymbolZset) {
		return 0, ErrCorruptedZsetScore
	}
	lKey := int(scoreKey[1])
	return ut.Bytes2Int64(scoreKey[lKey+3:])
}

func (z ZsetEncoder) DecodeMemberKey(memberKey []byte) ([]byte, error) {
	if len(memberKey) < 3 || memberKey[0] != (byte)(SymbolSet) {
		return nil, ErrCorruptedZsetMember
	}
	lKey := int(memberKey[1])
	return memberKey[lKey+2:], nil
}

func (z ZsetEncoder) QueueKey() []byte {
	return ut.ConcatBytearray(
		[]byte{(byte)(SymbolQueue)},
		z.key,
	)
}

func NewZsetEncoder(key []byte) (*ZsetEncoder, error) {
	if len(key) > MaxKeySize {
		return nil, ErrWrongNumOfArguments
	}
	return &ZsetEncoder{key: key}, nil
}
