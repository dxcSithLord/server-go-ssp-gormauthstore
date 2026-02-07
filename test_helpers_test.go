package gormauthstore

import (
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testIdentityBuilder provides a fluent API for constructing test identities.
// All fields default to safe, non-empty values so tests only need to set
// the fields relevant to their scenario.
type testIdentityBuilder struct {
	identity *ssp.SqrlIdentity
}

func newTestIdentity() *testIdentityBuilder {
	return &testIdentityBuilder{
		identity: &ssp.SqrlIdentity{
			Idk: "default-test-idk",
			Suk: "default-test-suk",
			Vuk: "default-test-vuk",
		},
	}
}

func (b *testIdentityBuilder) withIdk(idk string) *testIdentityBuilder {
	b.identity.Idk = idk
	return b
}

func (b *testIdentityBuilder) withSuk(suk string) *testIdentityBuilder {
	b.identity.Suk = suk
	return b
}

func (b *testIdentityBuilder) withVuk(vuk string) *testIdentityBuilder {
	b.identity.Vuk = vuk
	return b
}

func (b *testIdentityBuilder) withPidk(pidk string) *testIdentityBuilder {
	b.identity.Pidk = pidk
	return b
}

func (b *testIdentityBuilder) withSQRLOnly() *testIdentityBuilder {
	b.identity.SQRLOnly = true
	return b
}

func (b *testIdentityBuilder) withHardlock() *testIdentityBuilder {
	b.identity.Hardlock = true
	return b
}

func (b *testIdentityBuilder) withDisabled() *testIdentityBuilder {
	b.identity.Disabled = true
	return b
}

func (b *testIdentityBuilder) withBtn(btn int) *testIdentityBuilder {
	b.identity.Btn = btn
	return b
}

func (b *testIdentityBuilder) withRekeyed(rekeyed string) *testIdentityBuilder {
	b.identity.Rekeyed = rekeyed
	return b
}

func (b *testIdentityBuilder) build() *ssp.SqrlIdentity {
	return b.identity
}

// newTestStore creates an in-memory SQLite AuthStore ready for testing.
// The underlying connection pool is limited to 1 connection so that all
// goroutines share the same in-memory database (SQLite ":memory:" creates
// a separate database per connection).
func newTestStore(t *testing.T) *AuthStore {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return store
}

// newTestStoreWithDB creates an in-memory SQLite AuthStore and returns both
// the underlying *gorm.DB and the *AuthStore for tests that need direct DB access.
func newTestStoreWithDB(t *testing.T) (*gorm.DB, *AuthStore) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate failed: %v", err)
	}
	return db, store
}

// seedIdentity saves a test identity and fails the test on error.
func seedIdentity(t *testing.T, store *AuthStore, identity *ssp.SqrlIdentity) {
	t.Helper()
	if err := store.SaveIdentity(identity); err != nil {
		t.Fatalf("failed to seed identity %q: %v", identity.Idk, err)
	}
}
