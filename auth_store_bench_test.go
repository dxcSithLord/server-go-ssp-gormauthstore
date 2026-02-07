package gormauthstore

import (
	"fmt"
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// benchStore creates an in-memory SQLite AuthStore for benchmarks.
func benchStore(b *testing.B) *AuthStore {
	b.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatalf("failed to open database: %v", err)
	}
	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		b.Fatalf("AutoMigrate failed: %v", err)
	}
	return store
}

// PERF-001: Benchmark single-record FindIdentity.
func BenchmarkFindIdentity(b *testing.B) {
	store := benchStore(b)

	identity := &ssp.SqrlIdentity{
		Idk: "bench-find",
		Suk: "bench-suk",
		Vuk: "bench-vuk",
	}
	if err := store.SaveIdentity(identity); err != nil {
		b.Fatalf("seed failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.FindIdentity("bench-find")
	}
}

// PERF-002: Benchmark SaveIdentity (insert new records).
func BenchmarkSaveIdentity_Insert(b *testing.B) {
	store := benchStore(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		identity := &ssp.SqrlIdentity{
			Idk: fmt.Sprintf("bench-insert-%d", i),
			Suk: "suk",
			Vuk: "vuk",
		}
		_ = store.SaveIdentity(identity)
	}
}

// PERF-003: Benchmark SaveIdentity (update same record).
func BenchmarkSaveIdentity_Update(b *testing.B) {
	store := benchStore(b)

	identity := &ssp.SqrlIdentity{
		Idk: "bench-update",
		Suk: "suk",
		Vuk: "vuk",
	}
	if err := store.SaveIdentity(identity); err != nil {
		b.Fatalf("seed failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		identity.Btn = i % 4
		_ = store.SaveIdentity(identity)
	}
}

// PERF-004: Benchmark DeleteIdentity.
func BenchmarkDeleteIdentity(b *testing.B) {
	store := benchStore(b)

	// Pre-populate records to delete.
	for i := 0; i < b.N; i++ {
		identity := &ssp.SqrlIdentity{
			Idk: fmt.Sprintf("bench-del-%d", i),
			Suk: "suk",
			Vuk: "vuk",
		}
		if err := store.SaveIdentity(identity); err != nil {
			b.Fatalf("seed failed at %d: %v", i, err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.DeleteIdentity(fmt.Sprintf("bench-del-%d", i))
	}
}

// PERF-005: Benchmark ValidateIdk (pure CPU, no I/O).
func BenchmarkValidateIdk_AuthStore(b *testing.B) {
	validIdk := "k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateIdk(validIdk)
	}
}

// PERF-006: Benchmark concurrent FindIdentity reads.
func BenchmarkFindIdentity_Concurrent(b *testing.B) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		b.Fatalf("failed to open database: %v", err)
	}
	store := NewAuthStore(db)
	if err := store.AutoMigrate(); err != nil {
		b.Fatalf("AutoMigrate failed: %v", err)
	}

	identity := &ssp.SqrlIdentity{
		Idk: "bench-concurrent",
		Suk: "suk",
		Vuk: "vuk",
	}
	if err := store.SaveIdentity(identity); err != nil {
		b.Fatalf("seed failed: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = store.FindIdentity("bench-concurrent")
		}
	})
}
