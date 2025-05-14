package role

import (
	"simple-management-employee/internal/domain"

	"gorm.io/gorm"
)

type mysqlRoleRepository struct {
	db *gorm.DB
}

func NewMysqlRepository(db *gorm.DB) domain.RoleRepository {
	return &mysqlRoleRepository{db: db}
}

func (r *mysqlRoleRepository) Create(data *domain.Role) error {
	return r.db.Create(data).Error
}

func (r *mysqlRoleRepository) FindByID(id string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *mysqlRoleRepository) FindByName(name string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *mysqlRoleRepository) FindAll(page uint, size uint, filter *domain.Role) ([]*domain.Role, uint, error) {
	var roles []*domain.Role
	offset := (page - 1) * size
	query := r.db

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Offset(int(offset)).Limit(int(size)).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor uint
	if len(roles) > 0 {
		nextCursor = page + 1
	}

	return roles, nextCursor, nil

}

func (r *mysqlRoleRepository) Update(data *domain.Role) error {
	return r.db.Save(data).Error
}

func (r *mysqlRoleRepository) Delete(id string) error {
	return r.db.Delete(&domain.Role{}, id).Error
}

func (r *mysqlRoleRepository) Count(filter *domain.Role) (uint, error) {
	var count int64
	query := r.db.Model(&domain.Role{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return uint(count), nil
}

