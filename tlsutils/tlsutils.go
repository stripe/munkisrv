package tlsutils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/stripe/munkisrv/config"
)

// SetupTLSConfig creates a TLS configuration for the server based on the provided config
func SetupTLSConfig(cfg config.TLSConfig) (*tls.Config, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   getTLSVersion(cfg.MinVersion, tls.VersionTLS12),
		MaxVersion:   getTLSVersion(cfg.MaxVersion, tls.VersionTLS13),
	}

	// Setup client authentication if CA file is provided
	if cfg.CAFile != "" {
		caCert, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		tlsConfig.ClientCAs = caCertPool
		tlsConfig.ClientAuth = getClientAuthType(cfg.ClientAuth)
	}

	return tlsConfig, nil
}

// getTLSVersion converts a string version to tls.Version
func getTLSVersion(version string, defaultVersion uint16) uint16 {
	switch strings.ToLower(version) {
	case "1.0":
		return tls.VersionTLS10
	case "1.1":
		return tls.VersionTLS11
	case "1.2":
		return tls.VersionTLS12
	case "1.3":
		return tls.VersionTLS13
	default:
		return defaultVersion
	}
}

// getClientAuthType converts a string client auth type to tls.ClientAuthType
func getClientAuthType(clientAuth string) tls.ClientAuthType {
	switch strings.ToLower(clientAuth) {
	case "none":
		return tls.NoClientCert
	case "request":
		return tls.RequestClientCert
	case "require":
		return tls.RequireAnyClientCert
	case "verify-if-given":
		return tls.VerifyClientCertIfGiven
	case "require-and-verify":
		return tls.RequireAndVerifyClientCert
	default:
		return tls.NoClientCert
	}
}

// ValidateTLSConfig validates the TLS configuration
func ValidateTLSConfig(cfg config.TLSConfig) error {
	if !cfg.Enabled {
		return nil
	}

	// Check if certificate and key files exist
	if cfg.CertFile == "" {
		return fmt.Errorf("tls.cert_file is required when tls.enabled is true")
	}
	if cfg.KeyFile == "" {
		return fmt.Errorf("tls.key_file is required when tls.enabled is true")
	}

	// Check if files exist
	if _, err := os.Stat(cfg.CertFile); os.IsNotExist(err) {
		return fmt.Errorf("tls.cert_file does not exist: %s", cfg.CertFile)
	}
	if _, err := os.Stat(cfg.KeyFile); os.IsNotExist(err) {
		return fmt.Errorf("tls.key_file does not exist: %s", cfg.KeyFile)
	}

	// Validate client auth configuration
	if cfg.CAFile != "" {
		if _, err := os.Stat(cfg.CAFile); os.IsNotExist(err) {
			return fmt.Errorf("tls.ca_file does not exist: %s", cfg.CAFile)
		}
	}

	// Validate TLS versions
	if cfg.MinVersion != "" {
		if !isValidTLSVersion(cfg.MinVersion) {
			return fmt.Errorf("invalid tls.min_version: %s", cfg.MinVersion)
		}
	}
	if cfg.MaxVersion != "" {
		if !isValidTLSVersion(cfg.MaxVersion) {
			return fmt.Errorf("invalid tls.max_version: %s", cfg.MaxVersion)
		}
	}

	// Validate client auth type
	if cfg.ClientAuth != "" {
		if !isValidClientAuthType(cfg.ClientAuth) {
			return fmt.Errorf("invalid tls.client_auth: %s", cfg.ClientAuth)
		}
	}

	return nil
}

// isValidTLSVersion checks if a TLS version string is valid
func isValidTLSVersion(version string) bool {
	switch strings.ToLower(version) {
	case "1.0", "1.1", "1.2", "1.3":
		return true
	default:
		return false
	}
}

// isValidClientAuthType checks if a client auth type string is valid
func isValidClientAuthType(clientAuth string) bool {
	switch strings.ToLower(clientAuth) {
	case "none", "request", "require", "verify-if-given", "require-and-verify":
		return true
	default:
		return false
	}
}

// GetTLSInfo returns information about the TLS configuration
func GetTLSInfo(cfg config.TLSConfig) map[string]interface{} {
	info := map[string]interface{}{
		"enabled": cfg.Enabled,
	}

	if cfg.Enabled {
		info["cert_file"] = cfg.CertFile
		info["key_file"] = cfg.KeyFile
		info["ca_file"] = cfg.CAFile
		info["client_auth"] = cfg.ClientAuth
		info["min_version"] = cfg.MinVersion
		info["max_version"] = cfg.MaxVersion
	}

	return info
}
