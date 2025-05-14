package auth

import (
	"simple-management-employee/internal/domain"

	"gorm.io/gorm"
)

type mysqlAuthRepository struct {
	db *gorm.DB
}

func NewMysqlRepository(db *gorm.DB) domain.AuthRepository {
	return &mysqlAuthRepository{db: db}
}

func (r *mysqlAuthRepository) RegisterAdmin(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *mysqlAuthRepository) RegisterEmployee(user *domain.User) error {
	return r.db.Create(user).Error
}