package keyutils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// ParsePrivateKey parses a PEM encoded private key and returns a crypto.PrivateKey.
// It can be used for private keys passed in from environment variables or command line or files.
func ParsePrivateKey(privKeyPEM []byte, keyName string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(privKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode %s", keyName)
	}

	// The code below is based on tls.parsePrivateKey
	// https://cs.opensource.google/go/go/+/release-branch.go1.23:src/crypto/tls/tls.go;l=355-372
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("unmarshaled PKCS8 %s is not an RSA, ECDSA, or Ed25519 private key", keyName)
		}
	}
	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("failed to parse %s of type %s", keyName, block.Type)
}
