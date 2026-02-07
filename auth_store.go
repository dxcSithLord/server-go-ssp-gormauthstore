// Package gormauthstore implements the ssp.AuthStore interface for persisting
// SQRL authentication identities via the GORM ORM.
package gormauthstore

import (
	"errors"

	ssp "github.com/dxcSithLord/server-go-ssp"
	"gorm.io/gorm"
)

// identityRecord is a GORM v2 compatible model mirroring ssp.SqrlIdentity.
// The upstream SqrlIdentity struct uses legacy GORM v1 sql:"" tags (e.g.
// sql:"primary_key", sql:"-") that GORM v2 does not recognize. This model
// provides the correct GORM v2 tags while keeping the same database schema.
type identityRecord struct {
	Idk      string `gorm:"column:idk;primaryKey"`
	Suk      string `gorm:"column:suk"`
	Vuk      string `gorm:"column:vuk"`
	Pidk     string `gorm:"column:pidk"`
	SQRLOnly bool   `gorm:"column:sqrl_only"`
	Hardlock bool   `gorm:"column:hardlock"`
	Disabled bool   `gorm:"column:disabled"`
	Rekeyed  string `gorm:"column:rekeyed"`
	Btn      int    `gorm:"column:btn"`
}

// TableName returns the table name matching the GORM v1 convention for SqrlIdentity.
func (identityRecord) TableName() string {
	return "sqrl_identities"
}

// toRecord converts an ssp.SqrlIdentity to the GORM v2 model.
func toRecord(identity *ssp.SqrlIdentity) *identityRecord {
	return &identityRecord{
		Idk:      identity.Idk,
		Suk:      identity.Suk,
		Vuk:      identity.Vuk,
		Pidk:     identity.Pidk,
		SQRLOnly: identity.SQRLOnly,
		Hardlock: identity.Hardlock,
		Disabled: identity.Disabled,
		Rekeyed:  identity.Rekeyed,
		Btn:      identity.Btn,
	}
}

// toIdentity converts the GORM v2 model back to an ssp.SqrlIdentity.
func toIdentity(record *identityRecord) *ssp.SqrlIdentity {
	return &ssp.SqrlIdentity{
		Idk:      record.Idk,
		Suk:      record.Suk,
		Vuk:      record.Vuk,
		Pidk:     record.Pidk,
		SQRLOnly: record.SQRLOnly,
		Hardlock: record.Hardlock,
		Disabled: record.Disabled,
		Rekeyed:  record.Rekeyed,
		Btn:      record.Btn,
	}
}

// clearRecord wipes sensitive cryptographic fields from an identityRecord.
// Called after conversion to reduce the window where Suk/Vuk remain in memory.
func clearRecord(record *identityRecord) {
	WipeString(&record.Suk)
	WipeString(&record.Vuk)
}

// AuthStore is an ssp.AuthStore implementation using the gorm ORM.
type AuthStore struct {
	db *gorm.DB
}

// NewAuthStore creates an AuthStore using the passed in gorm instance.
func NewAuthStore(db *gorm.DB) *AuthStore {
	return &AuthStore{db}
}

// AutoMigrate uses gorm AutoMigrate to create/update the table holding the ssp.SqrlIdentity.
func (as *AuthStore) AutoMigrate() error {
	return as.db.AutoMigrate(&identityRecord{})
}

// FindIdentity implements ssp.AuthStore.
// Validates the idk before querying the database.
func (as *AuthStore) FindIdentity(idk string) (*ssp.SqrlIdentity, error) {
	if err := ValidateIdk(idk); err != nil {
		return nil, err
	}
	record := &identityRecord{}
	err := as.db.Where("idk = ?", idk).First(record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ssp.ErrNotFound
		}
		return nil, err
	}
	result := toIdentity(record)
	clearRecord(record)
	return result, nil
}

// SaveIdentity implements ssp.AuthStore.
// Validates the identity and its Idk before persisting.
func (as *AuthStore) SaveIdentity(identity *ssp.SqrlIdentity) error {
	if identity == nil {
		return ErrNilIdentity
	}
	if err := ValidateIdk(identity.Idk); err != nil {
		return err
	}
	record := toRecord(identity)
	err := as.db.Save(record).Error
	clearRecord(record)
	return err
}

// FindIdentitySecure retrieves a SQRL identity wrapped in a SecureIdentityWrapper.
// The wrapper provides RAII-style automatic cleanup of sensitive cryptographic
// material (Suk, Vuk) when Destroy() is called.
//
// Usage:
//
//	wrapper, err := store.FindIdentitySecure(idk)
//	if err != nil { return err }
//	defer wrapper.Destroy()
//	identity := wrapper.GetIdentity()
func (as *AuthStore) FindIdentitySecure(idk string) (*SecureIdentityWrapper, error) {
	identity, err := as.FindIdentity(idk)
	if err != nil {
		return nil, err
	}
	return NewSecureIdentityWrapper(identity), nil
}

// DeleteIdentity implements ssp.AuthStore.
// Validates the idk before executing the delete.
// Returns nil (no error) if the key does not exist.
func (as *AuthStore) DeleteIdentity(idk string) error {
	if err := ValidateIdk(idk); err != nil {
		return err
	}
	return as.db.Where("idk = ?", idk).Delete(&identityRecord{}).Error
}
