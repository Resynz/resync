package crypto

import (
	"encoding/base32"
)

const (
	_baseChars = "abcd2efgh3ijkl4mnop5qrst6uvwx7yz"
)

var (
	_base32Enc *base32.Encoding
)

func init() {
	_base32Enc = base32.NewEncoding(_baseChars).WithPadding(base32.NoPadding)
}

// B32Encode base32编码
func B32Encode(b []byte) string {
	return _base32Enc.EncodeToString(b)
}

// B32EncodeStr  string base32编码
func B32EncodeStr(s string) string {
	return _base32Enc.EncodeToString([]byte(s))
}

// B32Decode 解码
func B32Decode(s string) ([]byte, error) {
	return _base32Enc.DecodeString(s)
}
