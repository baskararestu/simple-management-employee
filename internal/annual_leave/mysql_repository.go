package annualleave

import (
	"simple-management-employee/internal/domain"

	"gorm.io/gorm"
)

type annualLeaveMysqlRepository struct {
	db *gorm.DB
}

func NewMysqlRepository(db *gorm.DB) domain.AnnualLeaveRepository {
	return &annualLeaveMysqlRepository{db: db}
}

func (m *annualLeaveMysqlRepository) Create(annualLeave *domain.AnnualLeave) error {
	return m.db.Create(annualLeave).Error
}

func (m *annualLeaveMysqlRepository) FindByID(id string) (*domain.AnnualLeave, error) {
	var annualLeave domain.AnnualLeave
	if err := m.db.Where("id = ?", id).First(&annualLeave).Error; err != nil {
		return nil, err
	}
	return &annualLeave, nil
}

func (r *annualLeaveMysqlRepository) FindAll(page uint, size uint, filter *domain.AnnualLeave) ([]*domain.AnnualLeave, uint, error) {
	var annualLeave []*domain.AnnualLeave
	offset := (page - 1) * size
	query := r.db

	if filter.User.FirstName != "" {
		query = query.Joins("JOIN users u ON u.id = annual_leaves.user_id").
			Where("u.first_name LIKE ?", "%"+filter.User.FirstName+"%")
	}

	if filter.User.LastName != "" {
		query = query.Joins("JOIN users u ON u.id = annual_leaves.user_id").
			Where("u.last_name LIKE ?", "%"+filter.User.LastName+"%")
	}

	if filter.User.Email != "" {
		query = query.Joins("JOIN users u ON u.id = annual_leaves.user_id").
			Where("u.email LIKE ?", "%"+filter.User.Email+"%")
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if !filter.StartDate.IsZero() {
		query = query.Where("start_date >= ?", filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		query = query.Where("end_date <= ?", filter.EndDate)
	}

	if err := query.Offset(int(offset)).Limit(int(size)).Find(&annualLeave).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor uint
	if len(annualLeave) > 0 {
		nextCursor = page + 1
	}

	return annualLeave, nextCursor, nil
}

func (m *annualLeaveMysqlRepository) Update(annualLeave *domain.AnnualLeave) error {
	return m.db.Save(annualLeave).Error
}

func (m *annualLeaveMysqlRepository) Delete(id string) error {
	return m.db.Where("id = ?", id).Delete(&domain.AnnualLeave{}).Error
}

func (m *annualLeaveMysqlRepository) FindByUserID(userID []string) ([]domain.AnnualLeave, error) {
	var annualLeaves []domain.AnnualLeave
	if err := m.db.Where("user_id IN ?", userID).Find(&annualLeaves).Error; err != nil {
		return nil, err
	}
	return annualLeaves, nil
}

func (m *annualLeaveMysqlRepository) FindByStatus(status string) ([]domain.AnnualLeave, error) {
	var annualLeaves []domain.AnnualLeave
	if err := m.db.Where("status = ?", status).Find(&annualLeaves).Error; err != nil {
		return nil, err
	}
	return annualLeaves, nil
}

func (m *annualLeaveMysqlRepository) FindByDateRange(startDate, endDate string) ([]domain.AnnualLeave, error) {
	var annualLeaves []domain.AnnualLeave
	if err := m.db.Where("start_date >= ? AND end_date <= ?", startDate, endDate).Find(&annualLeaves).Error; err != nil {
		return nil, err
	}
	return annualLeaves, nil
}

func (m *annualLeaveMysqlRepository) Approve(id string) error {
	return m.db.Model(&domain.AnnualLeave{}).Where("id = ?", id).Update("status", "approved").Error
}

func (m *annualLeaveMysqlRepository) Reject(id string) error {
	return m.db.Model(&domain.AnnualLeave{}).Where("id = ?", id).Update("status", "rejected").Error
}

func (m *annualLeaveMysqlRepository) Count(filter *domain.AnnualLeave) (uint, error) {
	var count int64
	query := m.db.Model(&domain.AnnualLeave{})

	if filter.User.FirstName != "" {
		query = query.Joins("JOIN users u ON u.id = annual_leaves.user_id").
			Where("u.first_name LIKE ?", "%"+filter.User.FirstName+"%")
	}

	if filter.User.LastName != "" {
		query = query.Joins("JOIN users u ON u.id = annual_leaves.user_id").
			Where("u.last_name LIKE ?", "%"+filter.User.LastName+"%")
	}

	if filter.User.Email != "" {
		query = query.Joins("JOIN users u ON u.id = annual_leaves.user_id").
			Where("u.email LIKE ?", "%"+filter.User.Email+"%")
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if !filter.StartDate.IsZero() {
		query = query.Where("start_date >= ?", filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		query = query.Where("end_date <= ?", filter.EndDate)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return uint(count), nil
}

func (m *annualLeaveMysqlRepository) GetAnnualLeavesInYear(userID string, year int) ([]*domain.AnnualLeave, error) {
	var annualLeaves []*domain.AnnualLeave
	if err := m.db.Where("user_id = ? AND YEAR(start_date) = ?", userID, year).Find(&annualLeaves).Error; err != nil {
		return nil, err
	}
	return annualLeaves, nil
}


