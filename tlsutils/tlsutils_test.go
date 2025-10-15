package tlsutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stripe/munkisrv/config"
)

func TestValidateTLSConfig(t *testing.T) {
	// Test valid TLS config (disabled)
	cfg := config.TLSConfig{Enabled: false}
	if err := ValidateTLSConfig(cfg); err != nil {
		t.Errorf("Expected no error for disabled TLS, got: %v", err)
	}

	// Test invalid TLS config (enabled but missing cert file)
	cfg = config.TLSConfig{Enabled: true}
	if err := ValidateTLSConfig(cfg); err == nil {
		t.Error("Expected error for missing cert_file")
	}

	// Test invalid TLS config (enabled but missing key file)
	cfg = config.TLSConfig{Enabled: true, CertFile: "/tmp/test.crt"}
	if err := ValidateTLSConfig(cfg); err == nil {
		t.Error("Expected error for missing key_file")
	}

	// Test invalid TLS version
	cfg = config.TLSConfig{Enabled: true, CertFile: "/tmp/test.crt", KeyFile: "/tmp/test.key", MinVersion: "1.4"}
	if err := ValidateTLSConfig(cfg); err == nil {
		t.Error("Expected error for invalid min_version")
	}

	// Test invalid client auth type
	cfg = config.TLSConfig{Enabled: true, CertFile: "/tmp/test.crt", KeyFile: "/tmp/test.key", ClientAuth: "invalid"}
	if err := ValidateTLSConfig(cfg); err == nil {
		t.Error("Expected error for invalid client_auth")
	}
}

func TestSetupTLSConfig(t *testing.T) {
	// Create temporary test certificates
	tempDir := t.TempDir()
	serverCert, serverKey := createTestCertificates(t, tempDir, "server")
	caCert := createTestCACertificate(t, tempDir)

	// Test TLS config with client authentication
	cfg := config.TLSConfig{
		Enabled:    true,
		CertFile:   serverCert,
		KeyFile:    serverKey,
		CAFile:     caCert,
		ClientAuth: "require-and-verify",
		MinVersion: "1.2",
		MaxVersion: "1.3",
	}

	tlsConfig, err := SetupTLSConfig(cfg)
	if err != nil {
		t.Fatalf("SetupTLSConfig failed: %v", err)
	}

	if tlsConfig == nil {
		t.Fatal("Expected TLS config to be created")
	}

	// Verify TLS configuration
	if len(tlsConfig.Certificates) != 1 {
		t.Error("Expected one certificate")
	}

	if tlsConfig.ClientAuth != tls.RequireAndVerifyClientCert {
		t.Errorf("Expected client auth type %v, got %v", tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)
	}

	if tlsConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("Expected min version %v, got %v", tls.VersionTLS12, tlsConfig.MinVersion)
	}

	if tlsConfig.MaxVersion != tls.VersionTLS13 {
		t.Errorf("Expected max version %v, got %v", tls.VersionTLS13, tlsConfig.MaxVersion)
	}

	// Test TLS config without client authentication
	cfg = config.TLSConfig{
		Enabled:    true,
		CertFile:   serverCert,
		KeyFile:    serverKey,
		ClientAuth: "none",
	}

	tlsConfig, err = SetupTLSConfig(cfg)
	if err != nil {
		t.Fatalf("SetupTLSConfig failed: %v", err)
	}

	if tlsConfig.ClientAuth != tls.NoClientCert {
		t.Errorf("Expected client auth type %v, got %v", tls.NoClientCert, tlsConfig.ClientAuth)
	}

	// Test disabled TLS
	cfg = config.TLSConfig{Enabled: false}
	tlsConfig, err = SetupTLSConfig(cfg)
	if err != nil {
		t.Fatalf("SetupTLSConfig failed: %v", err)
	}

	if tlsConfig != nil {
		t.Error("Expected nil TLS config when disabled")
	}
}

func TestGetTLSVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected uint16
	}{
		{"1.0", tls.VersionTLS10},
		{"1.1", tls.VersionTLS11},
		{"1.2", tls.VersionTLS12},
		{"1.3", tls.VersionTLS13},
		{"", tls.VersionTLS12},        // default
		{"invalid", tls.VersionTLS12}, // default
	}

	for _, test := range tests {
		result := getTLSVersion(test.input, tls.VersionTLS12)
		if result != test.expected {
			t.Errorf("getTLSVersion(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestGetClientAuthType(t *testing.T) {
	tests := []struct {
		input    string
		expected tls.ClientAuthType
	}{
		{"none", tls.NoClientCert},
		{"request", tls.RequestClientCert},
		{"require", tls.RequireAnyClientCert},
		{"verify-if-given", tls.VerifyClientCertIfGiven},
		{"require-and-verify", tls.RequireAndVerifyClientCert},
		{"", tls.NoClientCert},        // default
		{"invalid", tls.NoClientCert}, // default
	}

	for _, test := range tests {
		result := getClientAuthType(test.input)
		if result != test.expected {
			t.Errorf("getClientAuthType(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestGetTLSInfo(t *testing.T) {
	// Test disabled TLS
	cfg := config.TLSConfig{Enabled: false}
	info := GetTLSInfo(cfg)
	if info["enabled"] != false {
		t.Error("Expected enabled to be false")
	}

	// Test enabled TLS
	cfg = config.TLSConfig{
		Enabled:    true,
		CertFile:   "/path/to/cert.crt",
		KeyFile:    "/path/to/key.key",
		CAFile:     "/path/to/ca.crt",
		ClientAuth: "require-and-verify",
		MinVersion: "1.2",
		MaxVersion: "1.3",
	}
	info = GetTLSInfo(cfg)
	if info["enabled"] != true {
		t.Error("Expected enabled to be true")
	}
	if info["cert_file"] != "/path/to/cert.crt" {
		t.Error("Expected cert_file to match")
	}
	if info["key_file"] != "/path/to/key.key" {
		t.Error("Expected key_file to match")
	}
	if info["ca_file"] != "/path/to/ca.crt" {
		t.Error("Expected ca_file to match")
	}
	if info["client_auth"] != "require-and-verify" {
		t.Error("Expected client_auth to match")
	}
	if info["min_version"] != "1.2" {
		t.Error("Expected min_version to match")
	}
	if info["max_version"] != "1.3" {
		t.Error("Expected max_version to match")
	}
}

// Helper functions to create test certificates
func createTestCertificates(t *testing.T, tempDir, name string) (certFile, keyFile string) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Organization"},
			CommonName:   name + ".test",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{name + ".test", "localhost"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	// Write certificate to file
	certFile = filepath.Join(tempDir, name+".crt")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err := os.WriteFile(certFile, certPEM, 0644); err != nil {
		t.Fatalf("Failed to write certificate file: %v", err)
	}

	// Write private key to file
	keyFile = filepath.Join(tempDir, name+".key")
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err := os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	return certFile, keyFile
}

func createTestCACertificate(t *testing.T, tempDir string) string {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate CA private key: %v", err)
	}

	// Create CA certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA Organization"},
			CommonName:   "test-ca",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create CA certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create CA certificate: %v", err)
	}

	// Write CA certificate to file
	caFile := filepath.Join(tempDir, "ca.crt")
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err := os.WriteFile(caFile, caPEM, 0644); err != nil {
		t.Fatalf("Failed to write CA certificate file: %v", err)
	}

	return caFile
}
