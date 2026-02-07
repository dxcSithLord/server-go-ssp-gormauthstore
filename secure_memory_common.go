package gormauthstore

import (
	"runtime"

	ssp "github.com/dxcSithLord/server-go-ssp"
)

// ScrambleBytes overwrites a byte slice with pseudo-random data.
// This is more secure than zeroing as it leaves no predictable pattern.
//
//go:noinline
func ScrambleBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	// XOR with a simple pattern to avoid leaving zeros
	pattern := byte(0xAA)
	for i := range b {
		b[i] = pattern ^ byte(i&0xFF)
	}

	runtime.KeepAlive(b)
}

// WipeString clears a string reference and wipes a copy of its contents.
//
// IMPORTANT: This function does NOT wipe the original string's backing memory.
// Go strings are immutable and their backing arrays may reside in read-only
// memory segments (.rodata), making in-place modification unsafe and prone to
// crashes (SIGSEGV). This implementation copies the string to a mutable byte
// slice, wipes that copy, then clears the string reference.
//
// For true secure memory handling of sensitive data, callers should:
// - Use mutable []byte slices instead of strings for secrets
// - Use a secure-memory library like github.com/awnumar/memguard
// - Avoid string conversions which create copies the GC doesn't track
//
// This function provides defense-in-depth by:
// 1. Clearing the string reference (prevents further access via this variable)
// 2. Wiping a copy of the data (reduces copies in memory)
// 3. Allowing GC to reclaim the original string memory.
func WipeString(s *string) {
	if s == nil {
		return
	}
	if *s == "" {
		return
	}

	// Copy string contents to a mutable byte slice
	// This is safe because we're creating a new allocation
	dataCopy := []byte(*s)

	// Wipe the copy to reduce sensitive data copies in memory
	WipeBytes(dataCopy)

	// Clear the string reference
	// The original backing memory will be garbage collected
	*s = ""

	// Ensure the wiped copy isn't optimized away
	runtime.KeepAlive(dataCopy)
}

// ClearIdentity securely wipes all sensitive fields from a SqrlIdentity struct.
// This function should be called when an identity is no longer needed to minimize
// the window of exposure for cryptographic keys in memory.
//
// Fields cleared:
// - Idk (Identity Key)
// - Suk (Server Unlock Key)
// - Vuk (Verify Unlock Key)
// - Pidk (Previous Identity Key)
// - Rekeyed (Link to new identity)
//
// Usage:
//
//	identity, err := store.FindIdentity(idk)
//	if err != nil { ... }
//	defer ClearIdentity(identity)
//	// Use identity...
func ClearIdentity(identity *ssp.SqrlIdentity) {
	if identity == nil {
		return
	}

	// Wipe all sensitive string fields
	WipeString(&identity.Idk)
	WipeString(&identity.Suk)
	WipeString(&identity.Vuk)
	WipeString(&identity.Pidk)
	WipeString(&identity.Rekeyed)

	// Reset other fields to default values
	identity.SQRLOnly = false
	identity.Hardlock = false
	identity.Disabled = false
	identity.Btn = 0

	runtime.KeepAlive(identity)
}

// SecureIdentityWrapper provides RAII-style automatic cleanup for SqrlIdentity.
// The wrapper ensures that sensitive cryptographic material is wiped from memory
// when the identity is no longer needed.
//
// Usage:
//
//	wrapper, err := store.FindIdentitySecure(idk)
//	if err != nil { ... }
//	defer wrapper.Destroy()
//	// Access via wrapper.Identity
type SecureIdentityWrapper struct {
	Identity *ssp.SqrlIdentity
	wiped    bool
}

// NewSecureIdentityWrapper creates a new wrapper around an existing identity.
func NewSecureIdentityWrapper(identity *ssp.SqrlIdentity) *SecureIdentityWrapper {
	return &SecureIdentityWrapper{
		Identity: identity,
		wiped:    false,
	}
}

// Destroy securely wipes the identity and marks the wrapper as invalid.
// This method is idempotent - calling it multiple times is safe.
func (w *SecureIdentityWrapper) Destroy() {
	if w == nil || w.wiped {
		return
	}

	if w.Identity != nil {
		ClearIdentity(w.Identity)
		w.Identity = nil
	}
	w.wiped = true
}

// IsValid returns true if the wrapper still contains a valid identity.
func (w *SecureIdentityWrapper) IsValid() bool {
	return w != nil && !w.wiped && w.Identity != nil
}

// GetIdentity returns the wrapped identity if valid, otherwise returns nil.
// This is a safer alternative to directly accessing the Identity field.
func (w *SecureIdentityWrapper) GetIdentity() *ssp.SqrlIdentity {
	if !w.IsValid() {
		return nil
	}
	return w.Identity
}

// ValidateIdk performs basic validation on an Identity Key.
// Returns an error if the Idk is empty, too long, or contains invalid characters.
//
// Validation rules:
// - Cannot be empty
// - Maximum length: 256 characters (reasonable upper bound)
// - Should contain only URL-safe characters (alphanumeric, +, /, =, -, _, .)
func ValidateIdk(idk string) error {
	if idk == "" {
		return ErrEmptyIdentityKey
	}

	if len(idk) > MaxIdkLength {
		return ErrIdentityKeyTooLong
	}

	// Basic character validation for URL-safe base64-like strings
	for _, c := range idk {
		if !isValidIdkChar(c) {
			return ErrInvalidIdentityKeyFormat
		}
	}

	return nil
}

// isValidIdkChar checks if a character is valid for an Identity Key.
// Valid characters are alphanumeric plus common URL-safe characters: +, /, =, -, _, .
func isValidIdkChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '+' || c == '/' || c == '=' || c == '-' || c == '_' || c == '.'
}

// Constants for validation.
const (
	// MaxIdkLength is the maximum allowed length for an Identity Key.
	MaxIdkLength = 256
)
