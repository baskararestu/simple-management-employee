package user

import (
	"simple-management-employee/internal/domain"
	"simple-management-employee/pkg/xlogger"
)

type userService struct {
	repo domain.UserRepository
	annualLeaveSvc domain.AnnualLeaveService
}

func NewService(repo domain.UserRepository,annualLeaveSvc domain.AnnualLeaveService) domain.UserService {
	return &userService{repo: repo,annualLeaveSvc: annualLeaveSvc}
}

func (s *userService) FindByID(id string) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) FindByEmail(email string) (*domain.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *userService) FindAll(page uint, size uint, filter *domain.User) ([]*domain.UserResponseSvc, error) {
	userData, err := s.repo.FindAll(page, size, filter)
	if err != nil {
		xlogger.Logger.Error().Msgf("error finding users: %v", err)
		return nil, err
	}

	var userIDs []string
	for _, user := range userData {
		userIDs = append(userIDs, user.ID)
	}

	annualLeaveData, err := s.annualLeaveSvc.FindByUserID(userIDs)
	if err != nil {
		xlogger.Logger.Error().Msgf("error finding annual leaves: %v", err)
		return nil, err
	}

	annualLeaveMap := make(map[string]domain.AnnualLeaveResponse)
	for _, leave := range annualLeaveData {
		if _, exists := annualLeaveMap[leave.UserID]; !exists {
			annualLeaveMap[leave.UserID] = leave
		}
	}

	var responses []*domain.UserResponseSvc
	for _, user := range userData {
		var leaveResponse *domain.AnnualLeaveResponse
		if leave, ok := annualLeaveMap[user.ID]; ok {
			leaveResponse = &domain.AnnualLeaveResponse{
				ID: 	  leave.ID,
				UserID:    leave.UserID,
				StartDate: leave.StartDate,
				EndDate:   leave.EndDate,
				Status:    leave.Status,
				Reason:    leave.Reason,
			}
		}

		responses = append(responses, &domain.UserResponseSvc{
			ID:           user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email,
			Address:      user.Address,
			PhoneNumber:  user.PhoneNumber,
			Gender:       user.Gender,
			RoleName:     user.RoleName,
			AnnualLeave:  leaveResponse,
		})
	}

	return responses, nil
}

func (s *userService) Update(user *domain.User) error {
	return s.repo.Update(user)
}

func (s *userService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *userService) FindByRoleName(roleName string) ([]domain.User, error) {
	return s.repo.FindByRoleName(roleName)
}

func (s *userService) FindByRoleID(roleID string) ([]domain.User, error) {
	return s.repo.FindByRoleID(roleID)
}