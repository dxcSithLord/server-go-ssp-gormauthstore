//go:build integration

package gormauthstore

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestStore creates an in-memory SQLite AuthStore for testing.
func setupTestStore(t *testing.T) *AuthStore {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return store
}

// IT-001: Full Create-Read-Update-Read-Delete cycle.
func TestCRUDRoundTrip(t *testing.T) {
	store := setupTestStore(t)

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
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got: %v", err)
	}
}

// IT-002: FindIdentity returns ErrNotFound for non-existent keys.
func TestIntegration_FindIdentity_NotFound(t *testing.T) {
	store := setupTestStore(t)

	_, err := store.FindIdentity("nonexistent-key")
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}

// IT-003: Storing and retrieving multiple identities.
func TestSaveIdentity_MultipleRecords(t *testing.T) {
	store := setupTestStore(t)

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

	// Delete one and verify others remain.
	if err := store.DeleteIdentity("idk-bravo"); err != nil {
		t.Fatalf("DeleteIdentity failed: %v", err)
	}

	_, err := store.FindIdentity("idk-bravo")
	if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for deleted record, got: %v", err)
	}

	for _, idk := range []string{"idk-alpha", "idk-charlie"} {
		if _, err := store.FindIdentity(idk); err != nil {
			t.Errorf("FindIdentity(%q) should still exist: %v", idk, err)
		}
	}
}

// IT-004: Deleting a non-existent key is a no-op.
func TestIntegration_DeleteIdentity_NonExistent(t *testing.T) {
	store := setupTestStore(t)

	err := store.DeleteIdentity("nonexistent-key")
	if err != nil {
		t.Fatalf("DeleteIdentity on non-existent key should not error, got: %v", err)
	}
}

// IT-005: All SqrlIdentity fields round-trip correctly.
func TestIntegration_SaveIdentity_AllFields(t *testing.T) {
	store := setupTestStore(t)

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

// IT-006: Concurrent read/write operations are safe.
func TestIntegration_ConcurrentReadWrite(t *testing.T) {
	store := setupTestStore(t)

	// Seed a record for readers.
	identity := &ssp.SqrlIdentity{Idk: "concurrent-rw", Suk: "suk", Vuk: "vuk"}
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 20)

	// 10 concurrent readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := store.FindIdentity("concurrent-rw"); err != nil {
				errs <- fmt.Errorf("read: %w", err)
			}
		}()
	}

	// 10 concurrent writers to different keys
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			id := &ssp.SqrlIdentity{
				Idk: fmt.Sprintf("concurrent-w-%d", idx),
				Suk: "suk",
				Vuk: "vuk",
			}
			if err := store.SaveIdentity(id); err != nil {
				errs <- fmt.Errorf("write %d: %w", idx, err)
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent error: %v", err)
	}
}

// IT-007: Large dataset — insert and retrieve 1000 records.
func TestIntegration_LargeDataset(t *testing.T) {
	store := setupTestStore(t)

	const numIdentities = 1000
	for i := 0; i < numIdentities; i++ {
		identity := &ssp.SqrlIdentity{
			Idk: fmt.Sprintf("large-%04d", i),
			Suk: "suk",
			Vuk: "vuk",
		}
		if err := store.SaveIdentity(identity); err != nil {
			t.Fatalf("failed to save identity %d: %v", i, err)
		}
	}

	// Spot-check 100 records.
	for i := 0; i < 100; i++ {
		idk := fmt.Sprintf("large-%04d", i*10)
		if _, err := store.FindIdentity(idk); err != nil {
			t.Errorf("FindIdentity(%q) failed: %v", idk, err)
		}
	}
}

// IT-008: SQRL rekey workflow — Pidk links old and new identities.
func TestIntegration_RekeyWorkflow(t *testing.T) {
	store := setupTestStore(t)

	// Original identity.
	original := &ssp.SqrlIdentity{
		Idk: "original-idk",
		Suk: "original-suk",
		Vuk: "original-vuk",
	}
	if err := store.SaveIdentity(original); err != nil {
		t.Fatalf("save original failed: %v", err)
	}

	// New identity references old via Pidk.
	rekeyed := &ssp.SqrlIdentity{
		Idk:  "new-idk",
		Suk:  "new-suk",
		Vuk:  "new-vuk",
		Pidk: "original-idk",
	}
	if err := store.SaveIdentity(rekeyed); err != nil {
		t.Fatalf("save rekeyed failed: %v", err)
	}

	// Mark original as rekeyed.
	original.Rekeyed = "new-idk"
	if err := store.SaveIdentity(original); err != nil {
		t.Fatalf("update original failed: %v", err)
	}

	// Verify linkage.
	foundOriginal, err := store.FindIdentity("original-idk")
	if err != nil {
		t.Fatalf("find original failed: %v", err)
	}
	if foundOriginal.Rekeyed != "new-idk" {
		t.Errorf("Rekeyed: got %q, want %q", foundOriginal.Rekeyed, "new-idk")
	}

	foundNew, err := store.FindIdentity("new-idk")
	if err != nil {
		t.Fatalf("find rekeyed failed: %v", err)
	}
	if foundNew.Pidk != "original-idk" {
		t.Errorf("Pidk: got %q, want %q", foundNew.Pidk, "original-idk")
	}
}

// IT-009: FindIdentitySecure returns a wrapper and cleans up on Destroy.
func TestIntegration_FindIdentitySecure(t *testing.T) {
	store := setupTestStore(t)

	identity := &ssp.SqrlIdentity{
		Idk: "secure-integration",
		Suk: "secure-suk",
		Vuk: "secure-vuk",
	}
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	wrapper, err := store.FindIdentitySecure("secure-integration")
	if err != nil {
		t.Fatalf("FindIdentitySecure failed: %v", err)
	}
	if !wrapper.IsValid() {
		t.Fatal("wrapper should be valid")
	}

	got := wrapper.GetIdentity()
	if got.Suk != "secure-suk" {
		t.Errorf("Suk: got %q, want %q", got.Suk, "secure-suk")
	}

	wrapper.Destroy()
	if wrapper.IsValid() {
		t.Error("wrapper should be invalid after Destroy")
	}
}

// IT-010: Boolean field combinations persist correctly.
func TestIntegration_BooleanCombinations(t *testing.T) {
	store := setupTestStore(t)

	combos := []struct {
		idk      string
		sqrlOnly bool
		hardlock bool
		disabled bool
	}{
		{"bool-fff", false, false, false},
		{"bool-tff", true, false, false},
		{"bool-ftf", false, true, false},
		{"bool-fft", false, false, true},
		{"bool-ttt", true, true, true},
		{"bool-ttf", true, true, false},
		{"bool-ftt", false, true, true},
		{"bool-tft", true, false, true},
	}

	for _, c := range combos {
		identity := &ssp.SqrlIdentity{
			Idk:      c.idk,
			Suk:      "suk",
			Vuk:      "vuk",
			SQRLOnly: c.sqrlOnly,
			Hardlock: c.hardlock,
			Disabled: c.disabled,
		}
		if err := store.SaveIdentity(identity); err != nil {
			t.Fatalf("save %q failed: %v", c.idk, err)
		}
	}

	for _, c := range combos {
		found, err := store.FindIdentity(c.idk)
		if err != nil {
			t.Fatalf("find %q failed: %v", c.idk, err)
		}
		if found.SQRLOnly != c.sqrlOnly {
			t.Errorf("%s SQRLOnly: got %v, want %v", c.idk, found.SQRLOnly, c.sqrlOnly)
		}
		if found.Hardlock != c.hardlock {
			t.Errorf("%s Hardlock: got %v, want %v", c.idk, found.Hardlock, c.hardlock)
		}
		if found.Disabled != c.disabled {
			t.Errorf("%s Disabled: got %v, want %v", c.idk, found.Disabled, c.disabled)
		}
	}
}
