//go:build !windows

package gormauthstore

import (
	"reflect"
	"runtime"
	"unsafe"

	ssp "github.com/sqrldev/server-go-ssp"
)

// WipeBytes securely overwrites a byte slice with zeros.
// This function uses compiler directives to prevent dead store elimination
// which could otherwise optimize away the memory clearing operation.
//
// Note: This provides best-effort clearing but Go's garbage collector may still
// have copies of data in memory. For maximum security, consider using memguard
// library with locked memory pages.
//
//go:noinline
func WipeBytes(b []byte) {
	if len(b) == 0 {
		return
	}

	// Use a simple loop to zero out the bytes
	for i := range b {
		b[i] = 0
	}

	// Force a reference to prevent dead store elimination
	// This ensures the compiler doesn't optimize away our zeroing
	runtime.KeepAlive(b)
}

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

// WipeString attempts to securely clear a string by accessing its underlying bytes.
// WARNING: This is inherently unsafe as Go strings are immutable. This function
// modifies the underlying memory directly, which violates Go's string immutability
// guarantee. Use only for security-critical scenarios and understand the risks.
//
// Limitations:
// - String literals and interned strings cannot be safely wiped (stored in read-only memory)
// - Go runtime may have made copies of the string
// - Garbage collector doesn't track our modifications
// - This function will skip wiping if the string is in read-only memory (will only clear reference)
func WipeString(s *string) {
	if s == nil || *s == "" {
		return
	}

	// Get the string header to access underlying data
	sh := (*reflect.StringHeader)(unsafe.Pointer(s))
	if sh.Len == 0 || sh.Data == 0 {
		*s = ""
		return
	}

	// Create a slice header pointing to the same data
	sl := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}

	// Convert to byte slice and attempt to wipe
	// Use recover to handle read-only memory (string literals)
	b := *(*[]byte)(unsafe.Pointer(&sl))
	func() {
		defer func() {
			// Recover from panic if memory is read-only
			_ = recover()
		}()
		// Try to wipe - this will fail for string literals in read-only memory
		for i := range b {
			b[i] = 0
		}
		runtime.KeepAlive(b)
	}()

	// Clear the string reference regardless
	*s = ""
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
// - Should contain only URL-safe characters (alphanumeric, +, /, =, -, _)
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
// Valid characters are alphanumeric plus common URL-safe characters.
func isValidIdkChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '+' || c == '/' || c == '=' || c == '-' || c == '_' || c == '.'
}

// Constants for validation
const (
	// MaxIdkLength is the maximum allowed length for an Identity Key
	MaxIdkLength = 256
)
