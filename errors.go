package gormauthstore

import "errors"

// Package-specific errors for validation and security operations.
var (
	// ErrEmptyIdentityKey is returned when an empty identity key is provided.
	ErrEmptyIdentityKey = errors.New("identity key cannot be empty")

	// ErrIdentityKeyTooLong is returned when the identity key exceeds the maximum length.
	ErrIdentityKeyTooLong = errors.New("identity key exceeds maximum length of 256 characters")

	// ErrInvalidIdentityKeyFormat is returned when the identity key contains invalid characters.
	ErrInvalidIdentityKeyFormat = errors.New("identity key contains invalid characters")

	// ErrNilIdentity is returned when a nil identity is provided to an operation.
	ErrNilIdentity = errors.New("identity cannot be nil")

	// ErrNilDatabase is returned when the database connection is nil.
	ErrNilDatabase = errors.New("database connection cannot be nil")

	// ErrWrappedIdentityDestroyed is returned when accessing a destroyed wrapper.
	ErrWrappedIdentityDestroyed = errors.New("secure identity wrapper has been destroyed")
)
