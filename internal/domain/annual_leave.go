package domain

import (
	"time"
)

type AnnualLeave struct {
	ID         string    `gorm:"type:varchar(36);primaryKey"`
	UserID     string    `gorm:"type:varchar(36);not null"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
	Reason     string    `gorm:"type:varchar(255)"`
	Status string `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}


type AnnualLeaveRepository interface {
	Create(annualLeave *AnnualLeave) error
	FindByID(id string) (*AnnualLeave, error)
	FindAll(page uint, size uint, filter *AnnualLeave) ([]*AnnualLeave, uint, error)
	Update(annualLeave *AnnualLeave) error
	Delete(id string) error
	FindByUserID(userID []string) ([]AnnualLeave, error)
	FindByStatus(status string) ([]AnnualLeave, error)
	FindByDateRange(startDate, endDate string) ([]AnnualLeave, error)
	Approve(id string) error
	Reject(id string) error
	Count(filter *AnnualLeave) (uint, error)
	GetAnnualLeavesInYear(userID string, year int) ([]*AnnualLeave, error)
}

type AnnualLeaveService interface {
	Create(annualLeave *AnnualLeave) error
	FindByID(id string) (*AnnualLeave, error)
	FindAll(page uint, size uint, filter *AnnualLeave) ([]*AnnualLeaveResponse, uint, error)
	Update(annualLeave *AnnualLeave) error
	Delete(id string) error
	FindByUserID(userID []string) ([]AnnualLeaveResponse, error)
	FindByStatus(status string) ([]AnnualLeave, error)
	FindByDateRange(startDate, endDate string) ([]AnnualLeave, error)
	Approve(id string) error
	Reject(id string) error
	Count(filter *AnnualLeave) (uint, error)
	CountTotalLeaveDaysInYear(userID string, year int) (int, error)
}

type CreateAnnualLeaveRequest struct {
	UserID     string `json:"userId" validate:"required"`
	StartDate  string `json:"startDate" validate:"required"`
	EndDate    string `json:"endDate" validate:"required"`
	Reason     string `json:"reason" validate:"required"`
	Status     string `json:"status" validate:"required"`
}
	

type AnnualLeaveResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	StartDate  string    `json:"startDate"`
	EndDate    string    `json:"endDate"`
	Reason     string    `json:"reason"`
	Status     string    `json:"status"`
}
