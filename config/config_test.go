package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	configContent := `
server:
  host: "localhost"
  port: ":8080"
cloudfront:
  url: "https://test.cloudfront.net"
  key_id: "test-key-id"
  private_key: |
    -----BEGIN RSA PRIVATE KEY-----
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
    -----END RSA PRIVATE KEY-----
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Test loading the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify server config
	if config.Server.Host != "localhost" {
		t.Errorf("Expected server host 'localhost', got '%s'", config.Server.Host)
	}
	if config.Server.Port != ":8080" {
		t.Errorf("Expected server port ':8080', got '%s'", config.Server.Port)
	}

	// Verify cloudfront config
	if config.Cloudfront.URL != "https://test.cloudfront.net" {
		t.Errorf("Expected cloudfront URL 'https://test.cloudfront.net', got '%s'", config.Cloudfront.URL)
	}
	if config.Cloudfront.KeyID != "test-key-id" {
		t.Errorf("Expected cloudfront key ID 'test-key-id', got '%s'", config.Cloudfront.KeyID)
	}
	if config.Cloudfront.PrivateKey == "" {
		t.Error("Expected cloudfront private key to be set")
	}
}

func TestLoadConfigWithEnvOverrides(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	configContent := `
server:
  host: "localhost"
  port: ":8080"
cloudfront:
  url: "https://test.cloudfront.net"
  key_id: "test-key-id"
  private_key: "test-key"
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Set environment variables to override config
	os.Setenv("ENV_SERVER_HOST", "override-host")
	os.Setenv("ENV_SERVER_PORT", ":9090")
	defer os.Unsetenv("ENV_SERVER_HOST")
	defer os.Unsetenv("ENV_SERVER_PORT")

	// Test loading the config
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify environment overrides
	if config.Server.Host != "override-host" {
		t.Errorf("Expected server host 'override-host', got '%s'", config.Server.Host)
	}
	if config.Server.Port != ":9090" {
		t.Errorf("Expected server port ':9090', got '%s'", config.Server.Port)
	}
}

func TestLoadConfigInvalidPath(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent config file")
	}
}
