package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CompanyType string

const (
	CompanyTypeCorporation        CompanyType = "Corporations"
	CompanyTypeNonProfit          CompanyType = "Non Profit"
	CompanyTypeCooperative        CompanyType = "Cooperative"
	CompanyTypeSoleProprietorship CompanyType = "Sole Proprietorship"
)

var CompanyTypes = []CompanyType{
	CompanyTypeCorporation,
	CompanyTypeNonProfit,
	CompanyTypeCooperative,
	CompanyTypeSoleProprietorship,
}

type Company struct {
	ID                uuid.UUID   `gorm:"type:uuid;primary_key;" json:"ID,omitempty"`
	Name              string      `gorm:"type:varchar(50);unique;not null" json:"Name,omitempty" validate:"required"`
	Description       string      `gorm:"type:varchar(3000)" json:"Description,omitempty"`
	AmountOfEmployees int         `gorm:"not null" json:"AmountOfEmployees,omitempty" validate:"required"`
	Registered        bool        `gorm:"not null" json:"Registered"`
	Type              CompanyType `gorm:"type:varchar(20);not null" json:"Type,omitempty" validate:"required,company_type"`
	CreatedAt         time.Time   `json:"-"`
	UpdatedAt         time.Time   `json:"-"`
}

func (company *Company) BeforeCreate(tx *gorm.DB) (err error) {
	company.ID = uuid.New()
	return
}
