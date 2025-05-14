package domain

import (
	"time"
)

type User struct {
	ID          string    `gorm:"type:varchar(36);primaryKey"`
	FirstName   string    `gorm:"type:varchar(255)"`
	LastName    string    `gorm:"type:varchar(255)"`
	Email       string    `gorm:"type:varchar(255);unique"`
	Password    string    `gorm:"type:varchar(255)"`
	Address     *string    `gorm:"type:varchar(255)"`
	PhoneNumber *string    `gorm:"type:varchar(20)"`
	Gender 		*string   `gorm:"type:enum('MALE','FEMALE');default:null"` 
	RoleID      string    `gorm:"type:varchar(36);not null"`
	Role 	  Role     `gorm:"foreignKey:RoleID;references:ID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type UserRepository interface {
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll(page uint, size uint, filter *User) ([]*UserResponseCommon, error)
	Update(user *User) error
	Delete(id string) error
	FindByRoleName(roleName string) ([]User, error)
	FindByRoleID(roleID string) ([]User, error)
}

type UserService interface {
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll(page uint, size uint, filter *User) ([]*UserResponseSvc, error)
	Update(user *User) error
	Delete(id string) error
	FindByRoleName(roleName string) ([]User, error)
	FindByRoleID(roleID string) ([]User, error)
}

type CreateUserRequest struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Gender		*string `json:"gender"`
	RoleID      string `json:"roleId" validate:"required"`
}

type UpdateUserRequest struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Gender		*string `json:"gender"`
}

type UserResponseCommon struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Gender 	*string `json:"gender"`
	RoleName     string  `json:"roleName"`
}

type UserResponseSvc struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Gender 	*string `json:"gender"`
	RoleName     string  `json:"roleName"`
	AnnualLeave *AnnualLeaveResponse `json:"annualLeave"`
}