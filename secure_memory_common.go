package gormauthstore

import (
	"crypto/rand"
	"errors"
	"reflect"
	"unicode"
	"unsafe"

	ssp "github.com/sqrldev/server-go-ssp"
)

// MaxIdkLength defines the maximum length for an identity key
const MaxIdkLength = 256

// ScrambleBytes overwrites the byte slice with random data before wiping
func ScrambleBytes(b []byte) {
	if len(b) == 0 {
		return
	}
	// Overwrite with random data first
	_, _ = rand.Read(b)
	// Then wipe to zeros
	WipeBytes(b)
}

// WipeString securely wipes a string's underlying memory
// Note: This is best-effort as strings are immutable in Go
func WipeString(s *string) {
	if s == nil || len(*s) == 0 {
		return
	}
	// Convert string header to access underlying bytes
	// This is unsafe but necessary for secure wiping
	sh := (*reflect.StringHeader)(unsafe.Pointer(s))
	if sh.Data == 0 || sh.Len == 0 {
		return
	}
	// Create a byte slice pointing to the same memory
	b := unsafe.Slice((*byte)(unsafe.Pointer(sh.Data)), sh.Len)
	WipeBytes(b)
	*s = ""
}

// ClearIdentity securely wipes all sensitive fields of a SQRL identity
func ClearIdentity(identity *ssp.SqrlIdentity) {
	if identity == nil {
		return
	}
	WipeString(&identity.Idk)
	WipeString(&identity.Suk)
	WipeString(&identity.Vuk)
	WipeString(&identity.Pidk)
}

// SecureIdentityWrapper wraps a SqrlIdentity with secure memory handling
type SecureIdentityWrapper struct {
	identity *ssp.SqrlIdentity
	wiped    bool
}

// NewSecureIdentityWrapper creates a new secure wrapper for an identity
func NewSecureIdentityWrapper(identity *ssp.SqrlIdentity) *SecureIdentityWrapper {
	return &SecureIdentityWrapper{
		identity: identity,
		wiped:    false,
	}
}

// Identity returns the wrapped identity
func (w *SecureIdentityWrapper) Identity() *ssp.SqrlIdentity {
	if w.wiped {
		return nil
	}
	return w.identity
}

// Wipe securely clears the identity and marks it as wiped
func (w *SecureIdentityWrapper) Wipe() {
	if w.wiped || w.identity == nil {
		return
	}
	ClearIdentity(w.identity)
	w.wiped = true
}

// IsWiped returns true if the identity has been wiped
func (w *SecureIdentityWrapper) IsWiped() bool {
	return w.wiped
}

// ValidateIdk validates the format of an identity key
func ValidateIdk(idk string) error {
	if len(idk) == 0 {
		return errors.New("idk cannot be empty")
	}
	if len(idk) > MaxIdkLength {
		return errors.New("idk exceeds maximum length")
	}
	for _, c := range idk {
		if !isValidIdkChar(c) {
			return errors.New("idk contains invalid characters")
		}
	}
	return nil
}

// isValidIdkChar checks if a character is valid for an identity key
// Valid characters are alphanumeric, plus (+), slash (/), and equals (=)
func isValidIdkChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '+' || c == '/' || c == '=' || c == '-' || c == '_'
}
