package role

import (
	"simple-management-employee/internal/domain"
	"simple-management-employee/pkg/xlogger"
)

type roleService struct {
	repo domain.RoleRepository
}

func NewService(repo domain.RoleRepository) domain.RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) Create(req *domain.CreateRoleRequest) error {
	data := &domain.Role{
		Name:        req.Name,
	}
	return s.repo.Create(data)
}

func (s *roleService) FindByID(id string) (*domain.Role, error) {
	return s.repo.FindByID(id)
}

func (s *roleService) FindByName(name string) (*domain.Role, error) {
	return s.repo.FindByName(name)
}

func (s *roleService) FindAll(page uint, size uint, filter *domain.Role) ([]*domain.Role, uint,error){
	roles, nextCursor, err := s.repo.FindAll(page, size, filter)
	if err != nil {
		xlogger.Logger.Error().Msgf("error find all role: %v", err)
		return nil, 0, err
	}
	return roles, nextCursor, nil
}

func (s *roleService) Update(req *domain.Role) error {
	logger := xlogger.Logger
	role, err := s.repo.FindByID(req.ID)
	if err != nil {
		logger.Error().Msgf("error find role by id: %v", err)
		return err
	}
	return s.repo.Update(role)
}

func (s *roleService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *roleService) Count(filter *domain.Role) (uint, error) {
	count, err := s.repo.Count(filter)
	if err != nil {
		xlogger.Logger.Error().Msgf("error count role: %v", err)
		return 0, err
	}
	return count, nil
}