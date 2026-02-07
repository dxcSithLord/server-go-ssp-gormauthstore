package gormauthstore

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TC-001: NewAuthStore returns a valid, non-nil store.
func TestNewAuthStore_Success(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	store := NewAuthStore(db)
	if store == nil {
		t.Fatal("NewAuthStore returned nil")
	}
	if store.db == nil {
		t.Error("AuthStore.db is nil")
	}
}

// TC-002: AutoMigrate creates the sqrl_identities table.
func TestAutoMigrate_Success(t *testing.T) {
	db, store := newTestStoreWithDB(t)

	// Table should already exist (newTestStoreWithDB calls AutoMigrate).
	var count int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sqrl_identities'").Scan(&count)
	if count != 1 {
		t.Errorf("sqrl_identities table not created, count=%d", count)
	}

	// Extra explicit call to ensure the return value is nil on success.
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate returned error: %v", err)
	}
}

// TC-003: AutoMigrate is idempotent — calling it twice does not error.
func TestAutoMigrate_Idempotent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	store := NewAuthStore(db)

	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("first AutoMigrate failed: %v", err)
	}
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("second AutoMigrate failed: %v", err)
	}
}

// TC-004: SaveIdentity inserts a new record that can be read back.
func TestSaveIdentity_Insert(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().withIdk("tc004-insert").build()
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("SaveIdentity failed: %v", err)
	}

	found, err := store.FindIdentity("tc004-insert")
	if err != nil {
		t.Fatalf("FindIdentity failed: %v", err)
	}
	if found.Suk != identity.Suk {
		t.Errorf("Suk mismatch: got %q, want %q", found.Suk, identity.Suk)
	}
}

// TC-005: SaveIdentity updates an existing record (upsert behaviour).
func TestSaveIdentity_Update(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().withIdk("tc005-update").withSuk("original-suk").build()
	seedIdentity(t, store, identity)

	identity.Suk = "updated-suk"
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("SaveIdentity (update) failed: %v", err)
	}

	found, err := store.FindIdentity("tc005-update")
	if err != nil {
		t.Fatalf("FindIdentity failed: %v", err)
	}
	if found.Suk != "updated-suk" {
		t.Errorf("Suk not updated: got %q", found.Suk)
	}
}

// TC-006: SaveIdentity rejects nil identity.
func TestSaveIdentity_NilIdentity(t *testing.T) {
	store := newTestStore(t)

	err := store.SaveIdentity(nil)
	if err != ErrNilIdentity {
		t.Errorf("expected ErrNilIdentity, got %v", err)
	}
}

// TC-007: SaveIdentity rejects empty Idk.
func TestSaveIdentity_EmptyIdk(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().withIdk("").build()
	err := store.SaveIdentity(identity)
	if err != ErrEmptyIdentityKey {
		t.Errorf("expected ErrEmptyIdentityKey, got %v", err)
	}
}

// TC-008: FindIdentity returns a matching record.
func TestFindIdentity_Found(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().withIdk("tc008-find").build()
	seedIdentity(t, store, identity)

	found, err := store.FindIdentity("tc008-find")
	if err != nil {
		t.Fatalf("FindIdentity failed: %v", err)
	}
	if found == nil {
		t.Fatal("FindIdentity returned nil")
	}
	if found.Idk != "tc008-find" {
		t.Errorf("Idk mismatch: got %q", found.Idk)
	}
}

// TC-009: FindIdentity returns ssp.ErrNotFound for a missing key.
func TestFindIdentity_NotFound_Unit(t *testing.T) {
	store := newTestStore(t)

	_, err := store.FindIdentity("nonexistent-idk")
	if err != ssp.ErrNotFound {
		t.Errorf("expected ssp.ErrNotFound, got %v", err)
	}
}

// TC-010: FindIdentity rejects empty Idk.
func TestFindIdentity_EmptyIdk(t *testing.T) {
	store := newTestStore(t)

	_, err := store.FindIdentity("")
	if err != ErrEmptyIdentityKey {
		t.Errorf("expected ErrEmptyIdentityKey, got %v", err)
	}
}

// TC-011: FindIdentity rejects Idk exceeding max length.
func TestFindIdentity_IdkTooLong(t *testing.T) {
	store := newTestStore(t)

	longIdk := strings.Repeat("a", MaxIdkLength+1)
	_, err := store.FindIdentity(longIdk)
	if err != ErrIdentityKeyTooLong {
		t.Errorf("expected ErrIdentityKeyTooLong, got %v", err)
	}
}

// TC-012: FindIdentity rejects Idk with invalid characters.
func TestFindIdentity_InvalidCharacters(t *testing.T) {
	store := newTestStore(t)

	_, err := store.FindIdentity("invalid key with spaces")
	if err != ErrInvalidIdentityKeyFormat {
		t.Errorf("expected ErrInvalidIdentityKeyFormat, got %v", err)
	}
}

// TC-013: DeleteIdentity removes an existing record.
func TestDeleteIdentity_Exists(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().withIdk("tc013-delete").build()
	seedIdentity(t, store, identity)

	if err := store.DeleteIdentity("tc013-delete"); err != nil {
		t.Fatalf("DeleteIdentity failed: %v", err)
	}

	_, err := store.FindIdentity("tc013-delete")
	if err != ssp.ErrNotFound {
		t.Errorf("expected ssp.ErrNotFound after delete, got %v", err)
	}
}

// TC-014: DeleteIdentity on a non-existent key is a no-op.
func TestDeleteIdentity_NotExists(t *testing.T) {
	store := newTestStore(t)

	err := store.DeleteIdentity("nonexistent-idk")
	if err != nil {
		t.Errorf("DeleteIdentity should be a no-op for missing keys, got: %v", err)
	}
}

// TC-015: DeleteIdentity is idempotent — deleting twice does not error.
func TestDeleteIdentity_Idempotent(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().withIdk("tc015-idempotent").build()
	seedIdentity(t, store, identity)

	if err := store.DeleteIdentity("tc015-idempotent"); err != nil {
		t.Fatalf("first delete failed: %v", err)
	}
	if err := store.DeleteIdentity("tc015-idempotent"); err != nil {
		t.Fatalf("second delete should be idempotent, got: %v", err)
	}
}

// TC-016: DeleteIdentity rejects empty Idk.
func TestDeleteIdentity_EmptyIdk(t *testing.T) {
	store := newTestStore(t)

	err := store.DeleteIdentity("")
	if err != ErrEmptyIdentityKey {
		t.Errorf("expected ErrEmptyIdentityKey, got %v", err)
	}
}

// TC-017: All SqrlIdentity fields round-trip correctly.
func TestSaveIdentity_AllFields_Unit(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().
		withIdk("tc017-allfields").
		withSuk("test-suk-value").
		withVuk("test-vuk-value").
		withPidk("test-pidk-value").
		withSQRLOnly().
		withHardlock().
		withDisabled().
		withRekeyed("rekeyed-ref").
		withBtn(3).
		build()

	seedIdentity(t, store, identity)

	found, err := store.FindIdentity("tc017-allfields")
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

// TC-018: Concurrent FindIdentity reads are safe under the race detector.
func TestFindIdentity_Concurrent(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().withIdk("tc018-concurrent").build()
	seedIdentity(t, store, identity)

	const numReaders = 10
	var wg sync.WaitGroup
	errs := make(chan error, numReaders)

	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			_, err := store.FindIdentity("tc018-concurrent")
			if err != nil {
				errs <- err
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent read error: %v", err)
	}
}

// TC-019: Concurrent saves to different keys are safe.
func TestSaveIdentity_ConcurrentDifferentKeys(t *testing.T) {
	store := newTestStore(t)

	const numWriters = 10
	var wg sync.WaitGroup
	errs := make(chan error, numWriters)

	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(idx int) {
			defer wg.Done()
			id := newTestIdentity().withIdk(fmt.Sprintf("tc019-concurrent-%d", idx)).build()
			if err := store.SaveIdentity(id); err != nil {
				errs <- fmt.Errorf("writer %d: %w", idx, err)
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("concurrent save error: %v", err)
	}

	// Verify all records exist.
	for i := 0; i < numWriters; i++ {
		idk := fmt.Sprintf("tc019-concurrent-%d", i)
		if _, err := store.FindIdentity(idk); err != nil {
			t.Errorf("FindIdentity(%q) failed: %v", idk, err)
		}
	}
}

// TC-020: AuthStore implements the ssp.AuthStore interface (compile-time check).
func TestAuthStore_ImplementsInterface(t *testing.T) {
	var _ ssp.AuthStore = (*AuthStore)(nil)
}

// TC-021: SaveIdentity rejects Idk that is too long.
func TestSaveIdentity_IdkTooLong(t *testing.T) {
	store := newTestStore(t)

	longIdk := strings.Repeat("x", MaxIdkLength+1)
	identity := newTestIdentity().withIdk(longIdk).build()
	err := store.SaveIdentity(identity)
	if err != ErrIdentityKeyTooLong {
		t.Errorf("expected ErrIdentityKeyTooLong, got %v", err)
	}
}

// TC-022: SaveIdentity rejects Idk with invalid characters.
func TestSaveIdentity_InvalidIdk(t *testing.T) {
	store := newTestStore(t)

	identity := newTestIdentity().withIdk("invalid key!").build()
	err := store.SaveIdentity(identity)
	if err != ErrInvalidIdentityKeyFormat {
		t.Errorf("expected ErrInvalidIdentityKeyFormat, got %v", err)
	}
}

// TC-023: DeleteIdentity rejects Idk that is too long.
func TestDeleteIdentity_IdkTooLong(t *testing.T) {
	store := newTestStore(t)

	longIdk := strings.Repeat("x", MaxIdkLength+1)
	err := store.DeleteIdentity(longIdk)
	if err != ErrIdentityKeyTooLong {
		t.Errorf("expected ErrIdentityKeyTooLong, got %v", err)
	}
}

// TC-024: DeleteIdentity rejects Idk with invalid characters.
func TestDeleteIdentity_InvalidIdk(t *testing.T) {
	store := newTestStore(t)

	err := store.DeleteIdentity("bad key!")
	if err != ErrInvalidIdentityKeyFormat {
		t.Errorf("expected ErrInvalidIdentityKeyFormat, got %v", err)
	}
}

// TC-025: Multiple identities can be stored and retrieved independently.
func TestMultipleIdentities(t *testing.T) {
	store := newTestStore(t)

	ids := []string{"multi-alpha", "multi-bravo", "multi-charlie"}
	for _, idk := range ids {
		seedIdentity(t, store, newTestIdentity().withIdk(idk).withSuk("suk-"+idk).build())
	}

	for _, idk := range ids {
		found, err := store.FindIdentity(idk)
		if err != nil {
			t.Errorf("FindIdentity(%q) failed: %v", idk, err)
			continue
		}
		if found.Suk != "suk-"+idk {
			t.Errorf("Suk for %q: got %q, want %q", idk, found.Suk, "suk-"+idk)
		}
	}

	// Delete one, verify others survive.
	if err := store.DeleteIdentity("multi-bravo"); err != nil {
		t.Fatalf("DeleteIdentity failed: %v", err)
	}
	if _, err := store.FindIdentity("multi-bravo"); err != ssp.ErrNotFound {
		t.Errorf("expected ErrNotFound for deleted key, got %v", err)
	}
	for _, idk := range []string{"multi-alpha", "multi-charlie"} {
		if _, err := store.FindIdentity(idk); err != nil {
			t.Errorf("FindIdentity(%q) should still exist: %v", idk, err)
		}
	}
}

// TC-026: identityRecord.TableName returns the expected table name.
func TestIdentityRecord_TableName(t *testing.T) {
	r := identityRecord{}
	if r.TableName() != "sqrl_identities" {
		t.Errorf("TableName: got %q, want %q", r.TableName(), "sqrl_identities")
	}
}

// TC-027: toRecord and toIdentity are inverse operations.
func TestToRecord_ToIdentity_RoundTrip(t *testing.T) {
	original := &ssp.SqrlIdentity{
		Idk:      "roundtrip-idk",
		Suk:      "roundtrip-suk",
		Vuk:      "roundtrip-vuk",
		Pidk:     "roundtrip-pidk",
		SQRLOnly: true,
		Hardlock: true,
		Disabled: false,
		Rekeyed:  "roundtrip-rekeyed",
		Btn:      5,
	}

	record := toRecord(original)
	result := toIdentity(record)

	if result.Idk != original.Idk {
		t.Errorf("Idk: got %q, want %q", result.Idk, original.Idk)
	}
	if result.Suk != original.Suk {
		t.Errorf("Suk: got %q, want %q", result.Suk, original.Suk)
	}
	if result.Vuk != original.Vuk {
		t.Errorf("Vuk: got %q, want %q", result.Vuk, original.Vuk)
	}
	if result.Pidk != original.Pidk {
		t.Errorf("Pidk: got %q, want %q", result.Pidk, original.Pidk)
	}
	if result.SQRLOnly != original.SQRLOnly {
		t.Errorf("SQRLOnly: got %v, want %v", result.SQRLOnly, original.SQRLOnly)
	}
	if result.Hardlock != original.Hardlock {
		t.Errorf("Hardlock: got %v, want %v", result.Hardlock, original.Hardlock)
	}
	if result.Disabled != original.Disabled {
		t.Errorf("Disabled: got %v, want %v", result.Disabled, original.Disabled)
	}
	if result.Rekeyed != original.Rekeyed {
		t.Errorf("Rekeyed: got %q, want %q", result.Rekeyed, original.Rekeyed)
	}
	if result.Btn != original.Btn {
		t.Errorf("Btn: got %d, want %d", result.Btn, original.Btn)
	}
}
