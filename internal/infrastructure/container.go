package infrastructure

import (
	annualleave "simple-management-employee/internal/annual_leave"
	"simple-management-employee/internal/auth"
	"simple-management-employee/internal/config"
	"simple-management-employee/internal/domain"
	"simple-management-employee/internal/role"
	"simple-management-employee/internal/user"
	"simple-management-employee/pkg/xlogger"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

var (
	cfg config.Config
	
	userRepo domain.UserRepository
	authRepo domain.AuthRepository
	roleRepo domain.RoleRepository
	annualLeaveRepo domain.AnnualLeaveRepository

	userSvc domain.UserService
	authSvc domain.AuthService
	roleSvc domain.RoleService
	annualLeaveSvc domain.AnnualLeaveService
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	xlogger.Setup(cfg)
	xlogger.Logger.Info().Msgf("Config: %+v", cfg)
	dbSetup()

	userRepo = user.NewMysqlRepository(db)
	authRepo = auth.NewMysqlRepository(db)
	roleRepo = role.NewMysqlRepository(db)
	annualLeaveRepo = annualleave.NewMysqlRepository(db)

	roleSvc = role.NewService(roleRepo)
	annualLeaveSvc = annualleave.NewService(annualLeaveRepo)
	userSvc = user.NewService(userRepo,annualLeaveSvc)
	authSvc = auth.NewService(authRepo, userSvc,cfg.JwtSecret)
}
