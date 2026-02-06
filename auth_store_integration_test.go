//go:build integration

package gormauthstore

import (
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestCRUDRoundTrip exercises a full Create-Read-Update-Read-Delete cycle
// using an in-memory SQLite database to verify GORM v2 compatibility.
func TestCRUDRoundTrip(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	// --- Create ---
	identity := &ssp.SqrlIdentity{
		Idk:      "integration-test-idk",
		Suk:      "test-server-unlock-key",
		Vuk:      "test-verify-unlock-key",
		Pidk:     "",
		SQRLOnly: false,
		Hardlock: false,
		Disabled: false,
	}

	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("Create (SaveIdentity) failed: %v", err)
	}

	// --- Read ---
	found, err := store.FindIdentity("integration-test-idk")
	if err != nil {
		t.Fatalf("FindIdentity after create failed: %v", err)
	}
	if found.Idk != identity.Idk {
		t.Errorf("Idk mismatch: got %q, want %q", found.Idk, identity.Idk)
	}
	if found.Suk != identity.Suk {
		t.Errorf("Suk mismatch: got %q, want %q", found.Suk, identity.Suk)
	}
	if found.Vuk != identity.Vuk {
		t.Errorf("Vuk mismatch: got %q, want %q", found.Vuk, identity.Vuk)
	}

	// --- Update ---
	identity.Suk = "updated-server-unlock-key"
	identity.Disabled = true
	identity.Rekeyed = "new-idk-ref"

	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("Update (SaveIdentity) failed: %v", err)
	}

	// --- Read after update ---
	updated, err := store.FindIdentity("integration-test-idk")
	if err != nil {
		t.Fatalf("FindIdentity after update failed: %v", err)
	}
	if updated.Suk != "updated-server-unlock-key" {
		t.Errorf("Suk not updated: got %q", updated.Suk)
	}
	if !updated.Disabled {
		t.Error("Disabled not updated to true")
	}
	if updated.Rekeyed != "new-idk-ref" {
		t.Errorf("Rekeyed not updated: got %q", updated.Rekeyed)
	}

	// --- Delete ---
	if err := store.DeleteIdentity("integration-test-idk"); err != nil {
		t.Fatalf("DeleteIdentity failed: %v", err)
	}

	// --- Read after delete ---
	_, err = store.FindIdentity("integration-test-idk")
	if err != ssp.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got: %v", err)
	}
}

// TestFindIdentity_NotFound verifies ErrNotFound for non-existent keys.
func TestFindIdentity_NotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	_, err = store.FindIdentity("nonexistent-key")
	if err != ssp.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}

// TestSaveIdentity_MultipleRecords verifies storing and retrieving multiple identities.
func TestSaveIdentity_MultipleRecords(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	identities := []*ssp.SqrlIdentity{
		{Idk: "idk-alpha", Suk: "suk-alpha", Vuk: "vuk-alpha"},
		{Idk: "idk-bravo", Suk: "suk-bravo", Vuk: "vuk-bravo"},
		{Idk: "idk-charlie", Suk: "suk-charlie", Vuk: "vuk-charlie"},
	}

	for _, id := range identities {
		if err := store.SaveIdentity(id); err != nil {
			t.Fatalf("SaveIdentity(%q) failed: %v", id.Idk, err)
		}
	}

	for _, id := range identities {
		found, err := store.FindIdentity(id.Idk)
		if err != nil {
			t.Fatalf("FindIdentity(%q) failed: %v", id.Idk, err)
		}
		if found.Suk != id.Suk {
			t.Errorf("Suk mismatch for %q: got %q, want %q", id.Idk, found.Suk, id.Suk)
		}
	}

	// Delete one and verify others remain
	if err := store.DeleteIdentity("idk-bravo"); err != nil {
		t.Fatalf("DeleteIdentity failed: %v", err)
	}

	_, err = store.FindIdentity("idk-bravo")
	if err != ssp.ErrNotFound {
		t.Fatalf("expected ErrNotFound for deleted record, got: %v", err)
	}

	// Others should still exist
	for _, idk := range []string{"idk-alpha", "idk-charlie"} {
		if _, err := store.FindIdentity(idk); err != nil {
			t.Errorf("FindIdentity(%q) should still exist: %v", idk, err)
		}
	}
}

// TestDeleteIdentity_NonExistent verifies deleting a non-existent key does not error.
func TestDeleteIdentity_NonExistent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	// GORM performs a hard delete because identityRecord does not have a
	// gorm.DeletedAt field. Deleting a non-existent key is a no-op (no error).
	err = store.DeleteIdentity("nonexistent-key")
	if err != nil {
		t.Fatalf("DeleteIdentity on non-existent key should not error, got: %v", err)
	}
}

// TestSaveIdentity_AllFields verifies all SqrlIdentity fields round-trip correctly.
func TestSaveIdentity_AllFields(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}

	identity := &ssp.SqrlIdentity{
		Idk:      "full-field-test",
		Suk:      "suk-value",
		Vuk:      "vuk-value",
		Pidk:     "previous-idk",
		SQRLOnly: true,
		Hardlock: true,
		Disabled: true,
		Rekeyed:  "rekeyed-to-new",
		Btn:      3,
	}

	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("SaveIdentity failed: %v", err)
	}

	found, err := store.FindIdentity("full-field-test")
	if err != nil {
		t.Fatalf("FindIdentity failed: %v", err)
	}

	if found.Idk != identity.Idk {
		t.Errorf("Idk: got %q, want %q", found.Idk, identity.Idk)
	}
	if found.Suk != identity.Suk {
		t.Errorf("Suk: got %q, want %q", found.Suk, identity.Suk)
	}
	if found.Vuk != identity.Vuk {
		t.Errorf("Vuk: got %q, want %q", found.Vuk, identity.Vuk)
	}
	if found.Pidk != identity.Pidk {
		t.Errorf("Pidk: got %q, want %q", found.Pidk, identity.Pidk)
	}
	if found.SQRLOnly != identity.SQRLOnly {
		t.Errorf("SQRLOnly: got %v, want %v", found.SQRLOnly, identity.SQRLOnly)
	}
	if found.Hardlock != identity.Hardlock {
		t.Errorf("Hardlock: got %v, want %v", found.Hardlock, identity.Hardlock)
	}
	if found.Disabled != identity.Disabled {
		t.Errorf("Disabled: got %v, want %v", found.Disabled, identity.Disabled)
	}
	if found.Rekeyed != identity.Rekeyed {
		t.Errorf("Rekeyed: got %q, want %q", found.Rekeyed, identity.Rekeyed)
	}
	if found.Btn != identity.Btn {
		t.Errorf("Btn: got %d, want %d", found.Btn, identity.Btn)
	}
}
