package gormauthstore

import (
	"errors"
	"testing"

	ssp "github.com/dxcSithLord/server-go-ssp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// openTestDB creates an in-memory SQLite database for testing.
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory SQLite database: %v", err)
	}
	return db
}

func TestSave(t *testing.T) {
	db := openTestDB(t)
	gas := NewAuthStore(db)

	err := gas.AutoMigrate()
	if err != nil {
		t.Fatalf("couldn't automigrate to create sqrl_identity table: %v", err)
	}

	identity := &ssp.SqrlIdentity{
		Idk: "someidk",
		Suk: "server_unlock_key",
	}

	err = gas.SaveIdentity(identity)
	if err != nil {
		t.Fatalf("couldn't save identity: %v", err)
	}

	readback, err := gas.FindIdentity("someidk")
	if err != nil {
		t.Fatalf("couldn't find saved identity: %v", err)
	}

	if readback == nil || readback.Suk != "server_unlock_key" {
		t.Fatalf("readback identity not right: %#v", readback)
	}

	err = gas.DeleteIdentity("someidk")
	if err != nil {
		t.Fatalf("couldn't delete saved identity: %v", err)
	}

	_, err = gas.FindIdentity("someidk")
	if err == nil {
		t.Fatalf("should be deleted but isn't")
	} else if !errors.Is(err, ssp.ErrNotFound) {
		t.Fatalf("should be ErrNotFound but got: %v", err)
	}
}
