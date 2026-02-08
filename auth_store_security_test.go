package gormauthstore

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// openSecurityTestDB creates an in-memory SQLite database for security tests.
func openSecurityTestDB(t *testing.T) (*gorm.DB, *AuthStore) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db, store
}

// SEC-001: SQL injection prevention
// Verifies that common SQL injection payloads are rejected by input validation
// before reaching the database layer, and that the table remains intact.
func TestSQLInjectionPrevention(t *testing.T) {
	db, store := openSecurityTestDB(t)

	payloads := []string{
		"'; DROP TABLE sqrl_identities; --",
		"' OR '1'='1",
		"admin'--",
		"1' UNION SELECT * FROM users--",
		"'; DELETE FROM sqrl_identities WHERE '1'='1",
		"' OR 1=1--",
		"' UNION SELECT NULL, NULL, NULL--",
	}

	for _, payload := range payloads {
		t.Run("Find_"+payload, func(t *testing.T) {
			_, err := store.FindIdentity(payload)
			if err == nil {
				t.Errorf("SQL injection payload not rejected: %s", payload)
			}
		})

		t.Run("Delete_"+payload, func(t *testing.T) {
			err := store.DeleteIdentity(payload)
			if err == nil {
				t.Errorf("SQL injection payload not rejected on delete: %s", payload)
			}
		})

		t.Run("Save_"+payload, func(t *testing.T) {
			identity := &ssp.SqrlIdentity{
				Idk: payload,
				Suk: "test-suk",
				Vuk: "test-vuk",
			}
			err := store.SaveIdentity(identity)
			if err == nil {
				t.Errorf("SQL injection payload not rejected on save: %s", payload)
			}
		})
	}

	// Verify the table survived all injection attempts
	var count int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sqrl_identities'").Scan(&count)
	if count != 1 {
		t.Fatal("sqrl_identities table was dropped by an injection payload")
	}
}

// SEC-002: Length-based DoS prevention
// Verifies that extremely long identity keys are rejected before reaching
// the database, preventing resource exhaustion.
func TestDoSPrevention_LengthLimits(t *testing.T) {
	_, store := openSecurityTestDB(t)

	lengths := []int{257, 1000, 10000, 100000}
	for _, length := range lengths {
		t.Run(fmt.Sprintf("ValidateIdk/len_%d", length), func(t *testing.T) {
			longIdk := strings.Repeat("a", length)
			err := ValidateIdk(longIdk)
			if !errors.Is(err, ErrIdentityKeyTooLong) {
				t.Errorf("expected ErrIdentityKeyTooLong for length %d, got: %v", length, err)
			}
		})

		t.Run(fmt.Sprintf("FindIdentity/len_%d", length), func(t *testing.T) {
			longIdk := strings.Repeat("a", length)
			_, err := store.FindIdentity(longIdk)
			if !errors.Is(err, ErrIdentityKeyTooLong) {
				t.Errorf("expected ErrIdentityKeyTooLong for length %d, got: %v", length, err)
			}
		})

		t.Run(fmt.Sprintf("SaveIdentity/len_%d", length), func(t *testing.T) {
			longIdk := strings.Repeat("a", length)
			err := store.SaveIdentity(&ssp.SqrlIdentity{Idk: longIdk, Suk: "s", Vuk: "v"})
			if !errors.Is(err, ErrIdentityKeyTooLong) {
				t.Errorf("expected ErrIdentityKeyTooLong for length %d, got: %v", length, err)
			}
		})

		t.Run(fmt.Sprintf("DeleteIdentity/len_%d", length), func(t *testing.T) {
			longIdk := strings.Repeat("a", length)
			err := store.DeleteIdentity(longIdk)
			if !errors.Is(err, ErrIdentityKeyTooLong) {
				t.Errorf("expected ErrIdentityKeyTooLong for length %d, got: %v", length, err)
			}
		})
	}

	// Boundary: exactly at MaxIdkLength should be accepted
	exactIdk := strings.Repeat("a", MaxIdkLength)
	if err := ValidateIdk(exactIdk); err != nil {
		t.Errorf("Idk at exactly MaxIdkLength (%d) should be valid, got: %v", MaxIdkLength, err)
	}
}

// SEC-003: Character injection prevention (control characters)
// Verifies that newlines, tabs, null bytes, and other control characters
// are rejected to prevent log injection and header manipulation attacks.
func TestCharacterInjection_ControlChars(t *testing.T) {
	inputs := []struct {
		name  string
		value string
	}{
		{"newline", "idk\nwith\nnewlines"},
		{"carriage_return", "idk\rwith\rCR"},
		{"tab", "idk\twith\ttabs"},
		{"null_byte", "idk\x00with\x00nulls"},
		{"vertical_tab", "idk\vwith\vvtabs"},
		{"form_feed", "idk\fwith\fFF"},
		{"backspace", "idk\bwith\bBS"},
		{"space", "idk with spaces"},
		{"semicolon", "idk;semicolon"},
		{"backtick", "idk`backtick"},
		{"parentheses", "idk(parens)"},
		{"angle_brackets", "idk<angle>"},
		{"curly_braces", "idk{curly}"},
		{"pipe", "idk|pipe"},
		{"ampersand", "idk&amp"},
		{"dollar", "idk$dollar"},
		{"at_sign", "idk@at"},
		{"exclamation", "idk!bang"},
		{"hash", "idk#hash"},
		{"percent", "idk%pct"},
		{"asterisk", "idk*star"},
	}

	for _, tc := range inputs {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateIdk(tc.value)
			if !errors.Is(err, ErrInvalidIdentityKeyFormat) {
				t.Errorf("expected ErrInvalidIdentityKeyFormat for %q, got: %v", tc.name, err)
			}
		})
	}
}

// SEC-004: Sensitive data not in error messages
// Verifies that Suk, Vuk, and other cryptographic material never appear
// in error messages returned to callers.
func TestSensitiveDataNotInErrors(t *testing.T) {
	_, store := openSecurityTestDB(t)

	// Store a valid identity
	identity := &ssp.SqrlIdentity{
		Idk: "sensitive-test-idk",
		Suk: "SUPER_SECRET_SUK_VALUE_12345",
		Vuk: "SUPER_SECRET_VUK_VALUE_67890",
	}
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("setup SaveIdentity failed: %v", err)
	}

	// Trigger validation errors and check that secrets are not leaked
	errorCases := []struct {
		name string
		err  error
	}{
		{"empty_idk", func() error {
			return store.SaveIdentity(&ssp.SqrlIdentity{
				Idk: "",
				Suk: "SUPER_SECRET_SUK_VALUE_12345",
				Vuk: "SUPER_SECRET_VUK_VALUE_67890",
			})
		}()},
		{"nil_identity", store.SaveIdentity(nil)},
		{"invalid_format", func() error {
			_, err := store.FindIdentity("invalid key with spaces")
			return err
		}()},
		{"too_long", func() error {
			_, err := store.FindIdentity(strings.Repeat("a", MaxIdkLength+1))
			return err
		}()},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Fatal("expected an error but got nil")
			}
			errMsg := tc.err.Error()
			if strings.Contains(errMsg, "SUPER_SECRET_SUK_VALUE_12345") {
				t.Errorf("Suk leaked in error message: %s", errMsg)
			}
			if strings.Contains(errMsg, "SUPER_SECRET_VUK_VALUE_67890") {
				t.Errorf("Vuk leaked in error message: %s", errMsg)
			}
		})
	}
}

// SEC-005: Memory clearing verification
// Verifies that ClearIdentity wipes all sensitive fields from the struct.
func TestMemoryClearing(t *testing.T) {
	identity := &ssp.SqrlIdentity{
		Idk:      "memory-test-idk",
		Suk:      "secret-suk-value",
		Vuk:      "secret-vuk-value",
		Pidk:     "previous-idk-value",
		Rekeyed:  "rekeyed-value",
		SQRLOnly: true,
		Hardlock: true,
		Disabled: true,
		Btn:      3,
	}

	ClearIdentity(identity)

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
	if identity.SQRLOnly {
		t.Error("SQRLOnly not cleared")
	}
	if identity.Hardlock {
		t.Error("Hardlock not cleared")
	}
	if identity.Disabled {
		t.Error("Disabled not cleared")
	}
	if identity.Btn != 0 {
		t.Errorf("Btn not cleared: %d", identity.Btn)
	}
}

// SEC-005b: ClearIdentity is nil-safe.
func TestMemoryClearing_NilSafe(t *testing.T) {
	// Must not panic
	ClearIdentity(nil)
}

// SEC-006: Unicode normalisation attacks
// Verifies that zero-width characters, BiDi overrides, and other invisible
// Unicode code points are rejected, preventing homograph/normalisation attacks.
func TestUnicodeNormalizationAttacks(t *testing.T) {
	inputs := []struct {
		name  string
		value string
	}{
		{"zero_width_space", "idk\u200Binvisible"},
		{"right_to_left_override", "idk\u202Ertl"},
		{"byte_order_mark", "idk\uFEFFbom"},
		{"zero_width_joiner", "idk\u200Djoiner"},
		{"zero_width_non_joiner", "idk\u200Cnonjoiner"},
		{"left_to_right_mark", "idk\u200Eltr"},
		{"right_to_left_mark", "idk\u200Frtl"},
		{"soft_hyphen", "idk\u00ADhyphen"},
		{"non_breaking_space", "idk\u00A0nbsp"},
		{"homoglyph_cyrillic_a", "idk\u0430cyrillic"},
	}

	for _, tc := range inputs {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateIdk(tc.value)
			if !errors.Is(err, ErrInvalidIdentityKeyFormat) {
				t.Errorf("expected ErrInvalidIdentityKeyFormat for %q, got: %v", tc.name, err)
			}
		})
	}
}

// SEC-007: FindIdentitySecure wrapper behaviour
// Verifies that FindIdentitySecure returns a valid wrapper that provides
// access to the identity and properly cleans up on Destroy().
func TestFindIdentitySecure_Success(t *testing.T) {
	_, store := openSecurityTestDB(t)

	identity := &ssp.SqrlIdentity{
		Idk: "secure-find-test",
		Suk: "suk-secure",
		Vuk: "vuk-secure",
	}
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("SaveIdentity failed: %v", err)
	}

	wrapper, err := store.FindIdentitySecure("secure-find-test")
	if err != nil {
		t.Fatalf("FindIdentitySecure failed: %v", err)
	}

	if !wrapper.IsValid() {
		t.Fatal("wrapper should be valid before Destroy()")
	}

	got := wrapper.GetIdentity()
	if got == nil {
		t.Fatal("GetIdentity returned nil on valid wrapper")
	}
	if got.Idk != "secure-find-test" {
		t.Errorf("Idk mismatch: got %q, want %q", got.Idk, "secure-find-test")
	}
	if got.Suk != "suk-secure" {
		t.Errorf("Suk mismatch: got %q, want %q", got.Suk, "suk-secure")
	}

	// Destroy and verify cleanup
	wrapper.Destroy()

	if wrapper.IsValid() {
		t.Error("wrapper should be invalid after Destroy()")
	}
	if wrapper.GetIdentity() != nil {
		t.Error("GetIdentity should return nil after Destroy()")
	}
}

// SEC-008: FindIdentitySecure returns error for missing identity.
func TestFindIdentitySecure_NotFound(t *testing.T) {
	_, store := openSecurityTestDB(t)

	wrapper, err := store.FindIdentitySecure("nonexistent-idk")
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Errorf("expected ssp.ErrNotFound, got: %v", err)
	}
	if wrapper != nil {
		t.Error("expected nil wrapper on error")
	}
}

// SEC-009: FindIdentitySecure validates input.
func TestFindIdentitySecure_InvalidInput(t *testing.T) {
	_, store := openSecurityTestDB(t)

	cases := []struct {
		name string
		idk  string
		err  error
	}{
		{"empty", "", ErrEmptyIdentityKey},
		{"too_long", strings.Repeat("x", MaxIdkLength+1), ErrIdentityKeyTooLong},
		{"invalid_chars", "has spaces", ErrInvalidIdentityKeyFormat},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wrapper, err := store.FindIdentitySecure(tc.idk)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected %v, got: %v", tc.err, err)
			}
			if wrapper != nil {
				t.Error("expected nil wrapper on validation error")
			}
		})
	}
}

// SEC-010: SecureIdentityWrapper Destroy is idempotent.
func TestSecureIdentityWrapper_DestroyIdempotent(t *testing.T) {
	identity := &ssp.SqrlIdentity{
		Idk: "idempotent-test",
		Suk: "suk-value",
		Vuk: "vuk-value",
	}
	wrapper := NewSecureIdentityWrapper(identity)

	// Call Destroy multiple times -- should not panic.
	wrapper.Destroy()
	wrapper.Destroy()
	wrapper.Destroy()

	if wrapper.IsValid() {
		t.Error("wrapper should be invalid after Destroy()")
	}
}

// SEC-011: SecureIdentityWrapper nil safety.
func TestSecureIdentityWrapper_NilSafety(t *testing.T) {
	// Nil wrapper
	var wrapper *SecureIdentityWrapper
	if wrapper.IsValid() {
		t.Error("nil wrapper should not be valid")
	}
	if wrapper.GetIdentity() != nil {
		t.Error("nil wrapper GetIdentity should return nil")
	}
	// Destroy on nil should not panic
	wrapper.Destroy()
}

// SEC-012: clearRecord wipes sensitive fields from identityRecord.
func TestClearRecord_WipesSensitiveFields(t *testing.T) {
	record := &identityRecord{
		Idk: "test-idk",
		Suk: "secret-suk",
		Vuk: "secret-vuk",
	}

	clearRecord(record)

	if record.Suk != "" {
		t.Errorf("Suk not wiped: %q", record.Suk)
	}
	if record.Vuk != "" {
		t.Errorf("Vuk not wiped: %q", record.Vuk)
	}
	// Idk should be untouched (not sensitive in this context -- it's the lookup key)
	if record.Idk != "test-idk" {
		t.Errorf("Idk should be untouched: got %q", record.Idk)
	}
}

// SEC-013: Valid base64url-safe characters accepted.
func TestValidateIdk_AcceptsValidCharacters(t *testing.T) {
	// All allowed characters: alphanumeric + / = - _ .
	validIdks := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789",
		"k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E",
		"abc+def/ghi=jkl",
		"abc-def_ghi.jkl",
		"a",
	}

	for _, idk := range validIdks {
		if err := ValidateIdk(idk); err != nil {
			t.Errorf("ValidateIdk rejected valid idk %q: %v", idk, err)
		}
	}
}
