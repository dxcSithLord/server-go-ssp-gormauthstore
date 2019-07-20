package gormauthstore

import (
	"github.com/jinzhu/gorm"
	ssp "github.com/smw1218/sqrl-ssp"
)

// AuthStore is an ssp.AuthStore implementation using the gorm ORM
type AuthStore struct {
	db *gorm.DB
}

// NewAuthStore creates a AuthStore using the passed in gorm instance
func NewAuthStore(db *gorm.DB) *AuthStore {
	return &AuthStore{db}
}

// AutoMigrate uses the gorm.Automigrate to create/update the table holding the ssp.SqrlIdentity
func (as *AuthStore) AutoMigrate() error {
	return as.db.AutoMigrate(&ssp.SqrlIdentity{}).Error
}

// FindIdentity implements ssp.AuthStore
func (as *AuthStore) FindIdentity(idk string) (*ssp.SqrlIdentity, error) {
	identity := &ssp.SqrlIdentity{}
	err := as.db.Where("idk = ?", idk).First(identity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ssp.ErrNotFound
		}
		return nil, err
	}
	return identity, nil
}

// SaveIdentity implements ssp.AuthStore
func (as *AuthStore) SaveIdentity(identity *ssp.SqrlIdentity) error {
	return as.db.Save(identity).Error
}

// DeleteIdentity implements ssp.AuthStore
func (as *AuthStore) DeleteIdentity(idk string) error {
	return as.db.Where("idk = ?", idk).Delete(&ssp.SqrlIdentity{}).Error
}
