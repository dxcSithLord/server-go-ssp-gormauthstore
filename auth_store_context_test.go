package gormauthstore

import (
	"context"
	"errors"
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
)

// CTX-001: FindIdentityWithContext returns identity with valid context
func TestFindIdentityWithContext_ValidContext(t *testing.T) {
	store := newTestStore(t)
	identity := newTestIdentity().withIdk("ctx-find-idk").withSuk("ctx-suk").withVuk("ctx-vuk").build()
	seedIdentity(t, store, identity)

	ctx := context.Background()
	result, err := store.FindIdentityWithContext(ctx, "ctx-find-idk")
	if err != nil {
		t.Fatalf("FindIdentityWithContext failed: %v", err)
	}
	if result.Idk != "ctx-find-idk" {
		t.Fatalf("expected idk %q, got %q", "ctx-find-idk", result.Idk)
	}
	if result.Suk != "ctx-suk" {
		t.Fatalf("expected suk %q, got %q", "ctx-suk", result.Suk)
	}
}

// CTX-002: SaveIdentityWithContext persists with valid context
func TestSaveIdentityWithContext_ValidContext(t *testing.T) {
	store := newTestStore(t)
	identity := newTestIdentity().withIdk("ctx-save-idk").withSuk("save-suk").build()

	ctx := context.Background()
	err := store.SaveIdentityWithContext(ctx, identity)
	if err != nil {
		t.Fatalf("SaveIdentityWithContext failed: %v", err)
	}

	result, err := store.FindIdentity("ctx-save-idk")
	if err != nil {
		t.Fatalf("FindIdentity after save failed: %v", err)
	}
	if result.Suk != "save-suk" {
		t.Fatalf("expected suk %q, got %q", "save-suk", result.Suk)
	}
}

// CTX-003: DeleteIdentityWithContext deletes with valid context
func TestDeleteIdentityWithContext_ValidContext(t *testing.T) {
	store := newTestStore(t)
	identity := newTestIdentity().withIdk("ctx-del-idk").build()
	seedIdentity(t, store, identity)

	ctx := context.Background()
	err := store.DeleteIdentityWithContext(ctx, "ctx-del-idk")
	if err != nil {
		t.Fatalf("DeleteIdentityWithContext failed: %v", err)
	}

	_, err = store.FindIdentity("ctx-del-idk")
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got: %v", err)
	}
}

// CTX-004: FindIdentitySecureWithContext returns wrapper with valid context
func TestFindIdentitySecureWithContext_ValidContext(t *testing.T) {
	store := newTestStore(t)
	identity := newTestIdentity().withIdk("ctx-secure-idk").withSuk("sec-suk").withVuk("sec-vuk").build()
	seedIdentity(t, store, identity)

	ctx := context.Background()
	wrapper, err := store.FindIdentitySecureWithContext(ctx, "ctx-secure-idk")
	if err != nil {
		t.Fatalf("FindIdentitySecureWithContext failed: %v", err)
	}
	defer wrapper.Destroy()

	if !wrapper.IsValid() {
		t.Fatal("expected wrapper to be valid")
	}
	got := wrapper.GetIdentity()
	if got.Idk != "ctx-secure-idk" {
		t.Fatalf("expected idk %q, got %q", "ctx-secure-idk", got.Idk)
	}
}

// CTX-005: AutoMigrateWithContext succeeds with valid context
func TestAutoMigrateWithContext_ValidContext(t *testing.T) {
	// AutoMigrateWithContext is already tested implicitly by newTestStore,
	// but test it explicitly with a fresh DB.
	t.Helper()
	db := openTestDB(t)
	store := NewAuthStore(db)

	ctx := context.Background()
	err := store.AutoMigrateWithContext(ctx)
	if err != nil {
		t.Fatalf("AutoMigrateWithContext failed: %v", err)
	}

	// Verify table exists by saving an identity
	identity := newTestIdentity().withIdk("ctx-migrate-idk").build()
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("SaveIdentity after AutoMigrateWithContext failed: %v", err)
	}
}

// CTX-006: FindIdentityWithContext with cancelled context returns error
func TestFindIdentityWithContext_CancelledContext(t *testing.T) {
	store := newTestStore(t)
	identity := newTestIdentity().withIdk("ctx-cancel-find").build()
	seedIdentity(t, store, identity)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := store.FindIdentityWithContext(ctx, "ctx-cancel-find")
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}
}

// CTX-007: SaveIdentityWithContext with cancelled context returns error
func TestSaveIdentityWithContext_CancelledContext(t *testing.T) {
	store := newTestStore(t)
	identity := newTestIdentity().withIdk("ctx-cancel-save").build()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := store.SaveIdentityWithContext(ctx, identity)
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}
}

// CTX-008: DeleteIdentityWithContext with cancelled context returns error
func TestDeleteIdentityWithContext_CancelledContext(t *testing.T) {
	store := newTestStore(t)
	identity := newTestIdentity().withIdk("ctx-cancel-del").build()
	seedIdentity(t, store, identity)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := store.DeleteIdentityWithContext(ctx, "ctx-cancel-del")
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}
}

// CTX-009: WithContext methods validate idk same as originals
func TestWithContextMethods_ValidateIdk(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	// Empty idk
	_, err := store.FindIdentityWithContext(ctx, "")
	if !errors.Is(err, ErrEmptyIdentityKey) {
		t.Fatalf("FindIdentityWithContext: expected ErrEmptyIdentityKey, got %v", err)
	}

	err = store.DeleteIdentityWithContext(ctx, "")
	if !errors.Is(err, ErrEmptyIdentityKey) {
		t.Fatalf("DeleteIdentityWithContext: expected ErrEmptyIdentityKey, got %v", err)
	}

	// Nil identity
	err = store.SaveIdentityWithContext(ctx, nil)
	if !errors.Is(err, ErrNilIdentity) {
		t.Fatalf("SaveIdentityWithContext: expected ErrNilIdentity, got %v", err)
	}

	// Invalid characters
	_, err = store.FindIdentityWithContext(ctx, "bad<idk>")
	if !errors.Is(err, ErrInvalidIdentityKeyFormat) {
		t.Fatalf("FindIdentityWithContext: expected ErrInvalidIdentityKeyFormat, got %v", err)
	}
}

// CTX-010: FindIdentityWithContext not found returns ssp.ErrNotFound
func TestFindIdentityWithContext_NotFound(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	_, err := store.FindIdentityWithContext(ctx, "nonexistent-idk")
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("expected ssp.ErrNotFound, got %v", err)
	}
}

// CTX-011: Original methods still work (backward compatibility)
func TestOriginalMethods_BackwardCompatibility(t *testing.T) {
	store := newTestStore(t)

	// Save
	identity := newTestIdentity().withIdk("compat-idk").withSuk("compat-suk").build()
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("SaveIdentity failed: %v", err)
	}

	// Find
	result, err := store.FindIdentity("compat-idk")
	if err != nil {
		t.Fatalf("FindIdentity failed: %v", err)
	}
	if result.Suk != "compat-suk" {
		t.Fatalf("expected suk %q, got %q", "compat-suk", result.Suk)
	}

	// FindSecure
	wrapper, err := store.FindIdentitySecure("compat-idk")
	if err != nil {
		t.Fatalf("FindIdentitySecure failed: %v", err)
	}
	wrapper.Destroy()

	// Delete
	if err := store.DeleteIdentity("compat-idk"); err != nil {
		t.Fatalf("DeleteIdentity failed: %v", err)
	}

	_, err = store.FindIdentity("compat-idk")
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

// CTX-012: FindIdentitySecureWithContext with cancelled context returns error
func TestFindIdentitySecureWithContext_CancelledContext(t *testing.T) {
	store := newTestStore(t)
	identity := newTestIdentity().withIdk("ctx-cancel-secure").build()
	seedIdentity(t, store, identity)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := store.FindIdentitySecureWithContext(ctx, "ctx-cancel-secure")
	if err == nil {
		t.Fatal("expected error with cancelled context, got nil")
	}
}

// CTX-013: Full CRUD round-trip using WithContext methods
func TestWithContext_CRUDRoundTrip(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	// Create
	identity := newTestIdentity().
		withIdk("ctx-crud-idk").
		withSuk("crud-suk").
		withVuk("crud-vuk").
		withPidk("crud-pidk").
		withSQRLOnly().
		build()

	err := store.SaveIdentityWithContext(ctx, identity)
	if err != nil {
		t.Fatalf("SaveIdentityWithContext failed: %v", err)
	}

	// Read
	result, err := store.FindIdentityWithContext(ctx, "ctx-crud-idk")
	if err != nil {
		t.Fatalf("FindIdentityWithContext failed: %v", err)
	}
	if result.Suk != "crud-suk" || result.Vuk != "crud-vuk" {
		t.Fatalf("unexpected identity values: suk=%q vuk=%q", result.Suk, result.Vuk)
	}
	if !result.SQRLOnly {
		t.Fatal("expected SQRLOnly to be true")
	}

	// Update
	identity.Suk = "updated-suk"
	err = store.SaveIdentityWithContext(ctx, identity)
	if err != nil {
		t.Fatalf("SaveIdentityWithContext (update) failed: %v", err)
	}
	result, err = store.FindIdentityWithContext(ctx, "ctx-crud-idk")
	if err != nil {
		t.Fatalf("FindIdentityWithContext after update failed: %v", err)
	}
	if result.Suk != "updated-suk" {
		t.Fatalf("expected updated suk %q, got %q", "updated-suk", result.Suk)
	}

	// Delete
	err = store.DeleteIdentityWithContext(ctx, "ctx-crud-idk")
	if err != nil {
		t.Fatalf("DeleteIdentityWithContext failed: %v", err)
	}

	_, err = store.FindIdentityWithContext(ctx, "ctx-crud-idk")
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
