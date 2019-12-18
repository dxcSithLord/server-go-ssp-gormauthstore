package gormauthstore

// test expects a local postgres db with name "sqrl_test"
import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	ssp "github.com/sqrldev/server-go-ssp"
)

func TestSave(t *testing.T) {
	db, err := gorm.Open("postgres", "dbname=sqrl_test sslmode=disable")
	defer db.Close()
	if err != nil {
		t.Fatalf("failed connecting to sqrl_test db on a local postgres instance: %v", err)
	}
	db = db.Begin()
	defer db.Rollback()
	gas := NewAuthStore(db)

	err = gas.AutoMigrate()
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

	readback, err = gas.FindIdentity("someidk")
	if err == nil {
		t.Fatalf("should be deleted but isn't")
	} else {
		if err != ssp.ErrNotFound {
			t.Fatalf("should be ErrNotFound")
		}
	}

}
