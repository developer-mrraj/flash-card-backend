package repository

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AddressRepository interface {
	FindByUserID(userID uuid.UUID) ([]models.UserAddress, error)
	Create(address *models.UserAddress) error
	SetDefault(userID, addressID uuid.UUID) error
	Delete(userID, addressID uuid.UUID) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) FindByUserID(userID uuid.UUID) ([]models.UserAddress, error) {
	var addresses []models.UserAddress
	result := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses)
	return addresses, result.Error
}

func (r *addressRepository) Create(address *models.UserAddress) error {
	// If this is the first address, make it default
	var count int64
	r.db.Model(&models.UserAddress{}).Where("user_id = ?", address.UserID).Count(&count)
	if count == 0 {
		address.IsDefault = true
	}
	return r.db.Create(address).Error
}

func (r *addressRepository) SetDefault(userID, addressID uuid.UUID) error {
	// Clear all defaults for this user first
	if err := r.db.Model(&models.UserAddress{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		return err
	}
	// Set the new default
	return r.db.Model(&models.UserAddress{}).Where("id = ? AND user_id = ?", addressID, userID).Update("is_default", true).Error
}

func (r *addressRepository) Delete(userID, addressID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", addressID, userID).Delete(&models.UserAddress{}).Error
}
