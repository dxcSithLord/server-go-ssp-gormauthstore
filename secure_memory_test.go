package gormauthstore

import (
	"strings"
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
)

// TestWipeBytes verifies that byte slices are properly zeroed
func TestWipeBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{
			name:  "normal bytes",
			input: []byte("sensitive data here"),
		},
		{
			name:  "all ones",
			input: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:  "binary data",
			input: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		},
		{
			name:  "single byte",
			input: []byte{0xAB},
		},
		{
			name:  "large buffer",
			input: make([]byte, 4096),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill with known data
			for i := range tt.input {
				tt.input[i] = byte(i % 256)
			}

			WipeBytes(tt.input)

			// Verify all bytes are zero
			for i, b := range tt.input {
				if b != 0 {
					t.Errorf("byte at index %d is %x, expected 0", i, b)
				}
			}
		})
	}
}

func TestWipeBytes_EmptySlice(t *testing.T) {
	// Should not panic
	empty := []byte{}
	WipeBytes(empty)

	if len(empty) != 0 {
		t.Errorf("empty slice length changed: %d", len(empty))
	}
}

func TestWipeBytes_NilSlice(t *testing.T) {
	// Should not panic
	var nilSlice []byte
	WipeBytes(nilSlice)
}

func TestScrambleBytes(t *testing.T) {
	original := []byte("sensitive data")
	data := make([]byte, len(original))
	copy(data, original)

	ScrambleBytes(data)

	// Verify data has changed
	if string(data) == string(original) {
		t.Error("data was not scrambled")
	}

	// Verify data is not all zeros
	allZeros := true
	for _, b := range data {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Error("data was zeroed instead of scrambled")
	}
}

func TestScrambleBytes_EmptySlice(t *testing.T) {
	empty := []byte{}
	ScrambleBytes(empty)
	if len(empty) != 0 {
		t.Errorf("empty slice length changed: %d", len(empty))
	}
}

func TestWipeString(t *testing.T) {
	// Note: String literals in Go are stored in read-only memory and cannot be wiped.
	// WipeString only works on heap-allocated strings (e.g., from byte slices).
	// These tests use heap-allocated strings via byte slice conversion.

	tests := []struct {
		name  string
		input []byte // Use byte slices to ensure heap allocation
	}{
		{
			name:  "simple string",
			input: []byte("password123"),
		},
		{
			name:  "cryptographic key",
			input: []byte("aGVsbG8gd29ybGQ="),
		},
		{
			name:  "long string",
			input: []byte(strings.Repeat("secret", 100)),
		},
		{
			name:  "special characters",
			input: []byte("p@$$w0rd!@#$%^&*()"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create heap-allocated string from byte slice
			s := string(tt.input)
			WipeString(&s)

			if s != "" {
				t.Errorf("string not cleared: got %q, want empty", s)
			}
		})
	}
}

func TestWipeString_NilPointer(t *testing.T) {
	// Should not panic
	WipeString(nil)
}

func TestWipeString_EmptyString(t *testing.T) {
	s := ""
	WipeString(&s)
	if s != "" {
		t.Errorf("empty string changed: %q", s)
	}
}

func TestClearIdentity(t *testing.T) {
	// Use heap-allocated strings (from byte slices) to avoid read-only memory issues
	// In real usage, database reads return heap-allocated strings
	identity := &ssp.SqrlIdentity{
		Idk:      string([]byte("sensitive_idk_value")),
		Suk:      string([]byte("highly_sensitive_server_unlock_key")),
		Vuk:      string([]byte("verify_unlock_key_data")),
		Pidk:     string([]byte("previous_identity_key")),
		Rekeyed:  string([]byte("rekeyed_to_new_id")),
		SQRLOnly: true,
		Hardlock: true,
		Disabled: true,
		Btn:      42,
	}

	ClearIdentity(identity)

	// Verify all sensitive strings are cleared
	if identity.Idk != "" {
		t.Errorf("Idk not cleared: %q", identity.Idk)
	}
	if identity.Suk != "" {
		t.Errorf("Suk not cleared: %q", identity.Suk)
	}
	if identity.Vuk != "" {
		t.Errorf("Vuk not cleared: %q", identity.Vuk)
	}
	if identity.Pidk != "" {
		t.Errorf("Pidk not cleared: %q", identity.Pidk)
	}
	if identity.Rekeyed != "" {
		t.Errorf("Rekeyed not cleared: %q", identity.Rekeyed)
	}

	// Verify boolean fields are reset
	if identity.SQRLOnly {
		t.Error("SQRLOnly not reset to false")
	}
	if identity.Hardlock {
		t.Error("Hardlock not reset to false")
	}
	if identity.Disabled {
		t.Error("Disabled not reset to false")
	}

	// Verify integer field is reset
	if identity.Btn != 0 {
		t.Errorf("Btn not reset: got %d, want 0", identity.Btn)
	}
}

func TestClearIdentity_NilIdentity(t *testing.T) {
	// Should not panic
	ClearIdentity(nil)
}

func TestClearIdentity_EmptyFields(t *testing.T) {
	identity := &ssp.SqrlIdentity{
		Idk:      "",
		Suk:      "",
		Vuk:      "",
		Pidk:     "",
		Rekeyed:  "",
		SQRLOnly: false,
		Hardlock: false,
		Disabled: false,
		Btn:      0,
	}

	// Should not panic or change anything
	ClearIdentity(identity)

	if identity.Idk != "" || identity.Suk != "" {
		t.Error("empty fields should remain empty")
	}
}

func TestSecureIdentityWrapper_Basic(t *testing.T) {
	identity := &ssp.SqrlIdentity{
		Idk: string([]byte("test_idk")),
		Suk: string([]byte("test_suk")),
	}

	wrapper := NewSecureIdentityWrapper(identity)

	if !wrapper.IsValid() {
		t.Error("wrapper should be valid")
	}

	if wrapper.GetIdentity() != identity {
		t.Error("GetIdentity should return original identity")
	}

	wrapper.Destroy()

	if wrapper.IsValid() {
		t.Error("wrapper should be invalid after destroy")
	}

	if wrapper.GetIdentity() != nil {
		t.Error("GetIdentity should return nil after destroy")
	}
}

func TestSecureIdentityWrapper_DoubleDestroy(t *testing.T) {
	identity := &ssp.SqrlIdentity{
		Idk: string([]byte("test")),
	}

	wrapper := NewSecureIdentityWrapper(identity)

	// Should be idempotent
	wrapper.Destroy()
	wrapper.Destroy()
	wrapper.Destroy()

	if wrapper.IsValid() {
		t.Error("wrapper should remain invalid")
	}
}

func TestSecureIdentityWrapper_NilWrapper(t *testing.T) {
	var wrapper *SecureIdentityWrapper

	// Should not panic
	if wrapper.IsValid() {
		t.Error("nil wrapper should not be valid")
	}

	if wrapper.GetIdentity() != nil {
		t.Error("nil wrapper should return nil identity")
	}
}

func TestSecureIdentityWrapper_NilIdentity(t *testing.T) {
	wrapper := NewSecureIdentityWrapper(nil)

	if wrapper.IsValid() {
		t.Error("wrapper with nil identity should not be valid")
	}

	// Should not panic
	wrapper.Destroy()
}

func TestValidateIdk_Valid(t *testing.T) {
	validIdks := []string{
		"abc123",
		"ABCdef456",
		"a1b2c3d4e5f6",
		"idk-with-dashes",
		"idk_with_underscores",
		"idk.with.dots",
		"idk+with+plus",
		"idk/with/slashes",
		"base64encoded==",
		strings.Repeat("a", 256), // Maximum length
	}

	for _, idk := range validIdks {
		t.Run(idk[:min(20, len(idk))], func(t *testing.T) {
			err := ValidateIdk(idk)
			if err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
		})
	}
}

func TestValidateIdk_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		idk         string
		expectedErr error
	}{
		{
			name:        "empty string",
			idk:         "",
			expectedErr: ErrEmptyIdentityKey,
		},
		{
			name:        "too long",
			idk:         strings.Repeat("a", 257),
			expectedErr: ErrIdentityKeyTooLong,
		},
		{
			name:        "contains space",
			idk:         "idk with space",
			expectedErr: ErrInvalidIdentityKeyFormat,
		},
		{
			name:        "contains special chars",
			idk:         "idk@invalid",
			expectedErr: ErrInvalidIdentityKeyFormat,
		},
		{
			name:        "contains unicode",
			idk:         "idk-\u00e9",
			expectedErr: ErrInvalidIdentityKeyFormat,
		},
		{
			name:        "contains newline",
			idk:         "idk\nwith\nnewline",
			expectedErr: ErrInvalidIdentityKeyFormat,
		},
		{
			name:        "contains tab",
			idk:         "idk\twith\ttab",
			expectedErr: ErrInvalidIdentityKeyFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIdk(tt.idk)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if err != tt.expectedErr {
				t.Errorf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestIsValidIdkChar(t *testing.T) {
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/=-_."
	invalidChars := " !@#$%^&*()[]{}|\\:;\"'<>,?\n\t\r"

	for _, c := range validChars {
		if !isValidIdkChar(c) {
			t.Errorf("character %c should be valid", c)
		}
	}

	for _, c := range invalidChars {
		if isValidIdkChar(c) {
			t.Errorf("character %c should be invalid", c)
		}
	}
}

// Benchmark tests
func BenchmarkWipeBytes(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WipeBytes(data)
	}
}

func BenchmarkClearIdentity(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		identity := &ssp.SqrlIdentity{
			Idk:  "benchmark_idk_value",
			Suk:  "benchmark_suk_value",
			Vuk:  "benchmark_vuk_value",
			Pidk: "benchmark_pidk_value",
		}
		ClearIdentity(identity)
	}
}

func BenchmarkValidateIdk(b *testing.B) {
	idk := "valid_identity_key_123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateIdk(idk)
	}
}

func BenchmarkSecureWrapper(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		identity := &ssp.SqrlIdentity{
			Idk: "bench_idk",
			Suk: "bench_suk",
		}
		wrapper := NewSecureIdentityWrapper(identity)
		_ = wrapper.GetIdentity()
		wrapper.Destroy()
	}
}
