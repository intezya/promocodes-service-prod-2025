package pkg

import "github.com/intezya/pkglib"

type CryptoProvider struct{}

func (CryptoProvider) Salt() string {
	return string(pkglib.Crypto.Salt(32))
}

func (CryptoProvider) Encrypt(value string) string {
	return pkglib.Crypto.EncodeBase64(pkglib.Crypto.HashSHA256(value))
}
