package crypto

import (
	"encoding/base64"
)

const (
	_base64Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-."
)

var (
	_base64Enc *base64.Encoding
)

func init() {
	_base64Enc = base64.NewEncoding(_base64Chars).WithPadding(base64.NoPadding)
}

// Base64EncodeToString base64编码
func Base64EncodeToString(b []byte) string {
	return _base64Enc.EncodeToString(b)

}

// B64EncodeStr  string base64编码
func B64EncodeStr(s string) string {
	return _base64Enc.EncodeToString([]byte(s))
}

// Base64DecodeString base64解码
func Base64DecodeString(s string) ([]byte, error) {
	return _base64Enc.DecodeString(s)
}
