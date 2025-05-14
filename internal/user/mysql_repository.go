package user

import (
	"fmt"
	"simple-management-employee/internal/domain"

	"gorm.io/gorm"
)

type userMysqlRepository struct {
	db *gorm.DB
}

func NewMysqlRepository(db *gorm.DB) domain.UserRepository {
	return &userMysqlRepository{db: db}
}

func (m *userMysqlRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	if err := m.db.Where("id = ?", id).First(&user).Preload("role").Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *userMysqlRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := m.db.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}	
	return &user, nil
}

func (m *userMysqlRepository) FindAll(page uint, size uint, filter *domain.User) ([]*domain.UserResponseCommon, error) {
	var users []domain.User
	offset := (page - 1) * size
	query := m.db

	if filter.FirstName != "" {
		query = query.Where("first_name LIKE ?", "%"+filter.FirstName+"%")
	}

	if filter.LastName != "" {
		query = query.Where("last_name LIKE ?", "%"+filter.LastName+"%")
	}

	if filter.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.Email+"%")
	}

	query = query.Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})

	if err := query.Offset(int(offset)).Limit(int(size)).Find(&users).Error; err != nil {
		return nil, err
	}
	response := make([]*domain.UserResponseCommon, len(users))
	for i, user := range users {
		response[i] = &domain.UserResponseCommon{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			Address:     user.Address,
			PhoneNumber: user.PhoneNumber,
			Gender:      user.Gender,
			RoleName:    user.Role.Name, 
		}
	}

	return response, nil
}
func (m *userMysqlRepository) Update(user *domain.User) error {
	return m.db.Save(user).Error
}

func (m *userMysqlRepository) Delete(id string) error {
	return m.db.Where("id = ?", id).Delete(&domain.User{}).Error
}

func (m *userMysqlRepository) FindByRoleName(roleName string) ([]domain.User, error) {
	var users []domain.User
	if err := m.db.Joins("JOIN roles ON roles.id = users.role_id").Where("roles.name = ?", roleName).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (m *userMysqlRepository) FindByRoleID(roleID string) ([]domain.User, error) {
	var users []domain.User
	if err := m.db.Where("role_id = ?", roleID).Find(&users).Error; err != nil {
		return nil, err
	}
	fmt.Println(users, "users")
	return users, nil
}

