package domain

import (
	"time"
)

type Role struct {
	ID        string    `gorm:"type:varchar(36);primaryKey"`
	Name      string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}


type RoleRepository interface {
	Create(data *Role) error
	FindByID(id string) (*Role, error)
	FindByName(name string) (*Role, error)
	FindAll(page uint, size uint, filter *Role) ([]*Role, uint, error)
	Update(data *Role) error
	Delete(id string) error
	Count(filter *Role) (uint, error)
}

type RoleService interface {
	Create(req *CreateRoleRequest) error
	FindByID(id string) (*Role, error)
	FindByName(name string) (*Role, error)
	FindAll(page uint, size uint, filter *Role) ([]*Role, uint, error)
	Update(req *Role) error
	Delete(id string) error
	Count(filter *Role) (uint, error)
}

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}