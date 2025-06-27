package keyutils

import (
	"crypto/rsa"
	"testing"
)

func TestParsePrivateKeyRSA(t *testing.T) {
	// Test RSA private key
	rsaKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA8FwXbOgMCDuxjKASqKsBIcxklt1jAbn6QFWLUoDwE54nURlo
MZQ8yNSg0VsljegJudDhLy8IGuFTJCnjE3UQXuyAoOMk4lM6JMcFm71v5/rOuZiG
AfrtDWh931WHfAzXHztDxvrGVqTdzgiX6x+zeRXJ9B9jHaUcw7guWmCt2U8tVYTk
lDK+Gmfw9+/DmVSfc2BUbVmjbcSNpKqpGI+iQWElL65w4a516OG0Z+levQ1/uG7z
cfxMMbs8UjHQQ2KtlVRr1QTY4cPlr9/+umM6RPDsHyCFuJl00mtjz4zjkx0DPgAz
wiCiiRrbHv7Jew1FLrUEBJ0C0KyN7R+Yl61YzwIDAQABAoIBAAzkM5FwxKxwXy52
q2mGenIQn1iEGTpPej+XFvje14GF2v/7h94Y4EW5OcLgy5vX1SW1MU6xjBK9AROQ
d5Bkl/MvZhq69BB7fEPatM9MksLzbcEAkDds+OfeMdoXoUOjAKq5KAJ1Esw03Xye
c191/M9CvuksAcnmQCuzJjFMvCZKgTz4yNJJpi0f4uzRtJNeJ3s3cKDuG4IH+v/b
HCRF3qcXIxUTHbdNueTPubDIWtLVLQR3uHRbyIFs+Uu9WMjMplIp16SP16Mw9sbI
RpTah5ALJ3Gn4XYqrBkJPtB1d2Pis6EMz3vfQo6WXZYIR+BiXDctG3cD02Urq5rl
fbODMhkCgYEA+X5YRgZyZsdiPjwswRnS4MTb7xrOsVzdn3sCFmo9bcyaddGMvzUJ
kq3BDLxOQn48cQbFhVOs4Dj0MNadSrLApPaDXPMjWGR8ibdbCZ5PQToBSS3+LOOo
3Ng79KditwYOv1qbF8oENgfegVJSgWwbNB6nrOMoQmqZIzwq6E4bGaUCgYEA9qDE
odNqg0H5rXAwvKXohR5RbPFVPMe+lHJ7MxlhpoMQ9bTub7tlgrDsLibVChjMjYSr
6XzmaOXOfm40kMUvY6wh/LfSe/iu2qRSwCb+mmTUb2RBgo9zFkl5E7P6swsxHRZq
0vFt4+KwW4S51apsA7ITs+h7NCCEigo9QeAxVmMCgYAgkXCee3r1lbNqYlqJPoC7
nJcFKF+w4WmAxwLnwCiSq7HCDX+s+hRs1EeuDOq+XVIwguzH0btwbZ7avTk9JgZl
wlQ1jvufL0beh1PX9pVr81F1pw5V98X0RjnVXwBQ2faU3hP+z/0qvG48PW3NvTnz
3MiQlfqMaPPimJkVSBTbjQKBgQC2MdhxcDzMkM3BehMXGj2nMdmXcMW2bB13jwdC
naqNF2BNFAfdVQRNwyQHiDp0BhP/LBbQG6wfrD2bGxEMLg+vQ3esOaRuXy3VafWT
7HrEVl61l8vphs3PliGzE4/N+yOiSHBMO30iD9KXGXsrxIWdSU3S55k0zhz72Uqd
wuDP3wKBgQD09i4Xn/TjKDY1lJuUdp0HoswlpodeqXDWEIi3WEGPGJlYsu0EuCW+
Eh4UKDVUwJPr6yHZdW5FpSeXsuJLuxpAgCrgtZSr3j4N/L6hhUDrHOHU+zwUBWN+
Q6330VSxXxtQ47l5/pRwqs3Lc7K592FAvEpr99R650UY36J9F0471g==
-----END RSA PRIVATE KEY-----`

	key, err := ParsePrivateKey([]byte(rsaKey), "test RSA key")
	if err != nil {
		t.Fatalf("ParsePrivateKey failed: %v", err)
	}

	// Verify it's an RSA key
	rsaKeyType, ok := key.(*rsa.PrivateKey)
	if !ok {
		t.Fatalf("Expected RSA private key, got %T", key)
	}

	// Verify key has expected properties
	if rsaKeyType.N == nil {
		t.Error("RSA key modulus is nil")
	}
	if rsaKeyType.E == 0 {
		t.Error("RSA key exponent is zero")
	}
}

func TestParsePrivateKeyInvalidPEM(t *testing.T) {
	// Test with invalid PEM data
	invalidPEM := "not a valid PEM block"

	_, err := ParsePrivateKey([]byte(invalidPEM), "invalid key")
	if err == nil {
		t.Error("Expected error for invalid PEM data")
	}
}

func TestParsePrivateKeyEmptyData(t *testing.T) {
	// Test with empty data
	_, err := ParsePrivateKey([]byte(""), "empty key")
	if err == nil {
		t.Error("Expected error for empty key data")
	}
}

func TestParsePrivateKeyUnsupportedType(t *testing.T) {
	// Test with unsupported key type
	unsupportedKey := `-----BEGIN DSA PRIVATE KEY-----
MIIBuwIBAAKBgQDc+CZK9bBA9IU+gZUOc6FUGu7yO9WpNBZ9Z7CS30mHpXv3KtJX
-----END DSA PRIVATE KEY-----`

	_, err := ParsePrivateKey([]byte(unsupportedKey), "unsupported key")
	if err == nil {
		t.Error("Expected error for unsupported key type")
	}
}
