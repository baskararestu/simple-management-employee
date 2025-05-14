package domain

type AuthService interface {
	Login(req *LoginRequest) (*LoginResponse, error)
	RegisterAdmin(req *User) error
	RegisterEmployee(req *User) error
}

type AuthRepository interface {
	RegisterAdmin(user *User) error
	RegisterEmployee(user *User) error
}

type RegisterAdminRequest struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	RoleID      string `json:"roleId" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token *string `json:"token"`
}

type RegisterEmployeeRequest struct {
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Gender		string `json:"gender" validate:"required"`
	RoleID      string `json:"roleId" validate:"required"`
}