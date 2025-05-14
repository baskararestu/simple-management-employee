package auth

import (
	"simple-management-employee/internal/domain"
	"simple-management-employee/internal/utilities"
	"simple-management-employee/pkg/xlogger"
)

type authService struct{
	repo domain.AuthRepository
	userSvc domain.UserService
	jwtSecret string
}

func NewService(repo domain.AuthRepository, userSvc domain.UserService,	jwtSecret string) domain.AuthService {
	return &authService{repo: repo,userSvc: userSvc, jwtSecret: jwtSecret}
}

func (s *authService) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.userSvc.FindByEmail(req.Email)
	if err != nil {
		xlogger.Logger.Error().Msgf("error finding user: %v", err)
		return nil, err
	}

	if err := utilities.CheckPasswordHash(req.Password, user.Password); err != nil {
		xlogger.Logger.Error().Msgf("error checking password: %v", err)
		return nil, err
	}
	token, err := utilities.GenerateToken(user.ID, user.RoleID,s.jwtSecret)
	if err != nil {
		xlogger.Logger.Error().Msgf("error generating token: %v", err)
		return nil, err
	}
	
	return &domain.LoginResponse{
		Token: token,
	}, nil
}

func (s *authService) RegisterAdmin(req *domain.User) error {
	hashedPassword, err := utilities.HashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashedPassword
	return s.repo.RegisterAdmin(req)
}

func (s *authService) RegisterEmployee(req *domain.User) error {
	hashedPassword, err := utilities.HashPassword(req.Password)
	if err != nil {
		return err
	}
	req.Password = hashedPassword
	return s.repo.RegisterEmployee(req)
}
