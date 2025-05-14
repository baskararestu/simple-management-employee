package annualleave

import (
	"fmt"
	"simple-management-employee/internal/domain"
)

type annualLeaveService struct {
	repo domain.AnnualLeaveRepository
}

func NewService(repo domain.AnnualLeaveRepository) domain.AnnualLeaveService {
	return &annualLeaveService{repo: repo}
}

func (s *annualLeaveService) Create(annualLeave *domain.AnnualLeave) error {
	const maxAnnualLeaveDays = 5

	requestedDays := int(annualLeave.EndDate.Sub(annualLeave.StartDate).Hours()/24) + 1
	if requestedDays <= 0 {
		return fmt.Errorf("invalid leave duration")
	}

	year := annualLeave.StartDate.Year()

	usedDays, err := s.CountTotalLeaveDaysInYear(annualLeave.UserID, year)
	if err != nil {
		return err
	}

	if usedDays+requestedDays > maxAnnualLeaveDays {
		return domain.NewError(400, "exceeds maximum annual leave days")
	}

	return s.repo.Create(annualLeave)
}


func (s *annualLeaveService) FindByID(id string) (*domain.AnnualLeave, error) {
	return s.repo.FindByID(id)
}

func (s *annualLeaveService) FindAll(page uint, size uint, filter *domain.AnnualLeave) ([]*domain.AnnualLeaveResponse, uint,error) {
	annualLeaves, nextPage, err := s.repo.FindAll(page, size, filter)
	if err != nil {
		return nil,0, err
	}
	var annualLeaveResponses []*domain.AnnualLeaveResponse
	for _, leave := range annualLeaves {
		annualLeaveResponses = append(annualLeaveResponses, &domain.AnnualLeaveResponse{
			ID:         leave.ID,
			UserID:     leave.UserID,
			StartDate:  leave.StartDate.Format("2006-01-02"),
			EndDate:    leave.EndDate.Format("2006-01-02"),
			Status:     leave.Status,
			Reason:     leave.Reason,
		})
	}
	return annualLeaveResponses, nextPage, nil
}

func (s *annualLeaveService) Update(annualLeave *domain.AnnualLeave) error {
	return s.repo.Update(annualLeave)
}

func (s *annualLeaveService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *annualLeaveService) FindByUserID(userID []string) ([]domain.AnnualLeaveResponse, error) {
	data,err:= s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var annualLeaveResponses []domain.AnnualLeaveResponse
	for _, leave := range data {
		annualLeaveResponses = append(annualLeaveResponses, domain.AnnualLeaveResponse{
			ID:         leave.ID,
			UserID:     leave.UserID,
			StartDate:  leave.StartDate.Format("2006-01-02"),
			EndDate:    leave.EndDate.Format("2006-01-02"),
			Status:     leave.Status,
			Reason:     leave.Reason,
		})
	}
	return annualLeaveResponses, nil
}

func (s *annualLeaveService) FindByStatus(status string) ([]domain.AnnualLeave, error) {
	return s.repo.FindByStatus(status)
}

func (s *annualLeaveService) FindByDateRange(startDate, endDate string) ([]domain.AnnualLeave, error) {
	return s.repo.FindByDateRange(startDate, endDate)
}

func (s *annualLeaveService) Approve(id string) error {
	return s.repo.Approve(id)
}

func (s *annualLeaveService) Reject(id string) error {
	return s.repo.Reject(id)
}

func (s *annualLeaveService) Count(filter *domain.AnnualLeave) (uint, error) {
	return s.repo.Count(filter)
}

func (s *annualLeaveService) CountTotalLeaveDaysInYear(userID string, year int) (int, error) {
	leaves, err := s.repo.GetAnnualLeavesInYear(userID, year)
	if err != nil {
		return 0, err
	}

	totalDays := 0
	for _, leave := range leaves {
		days := int(leave.EndDate.Sub(leave.StartDate).Hours()/24) + 1
		totalDays += days
	}
	return totalDays, nil
}
