package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/vcsfrl/xm/internal/model"
	"gorm.io/gorm"
)

var ErrCompanyService = errors.New("company service error")

type Company struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewCompanyService(db *gorm.DB, validator *validator.Validate) *Company {
	return &Company{db: db, validator: validator}
}

func (s *Company) Create(company *model.Company) error {

	// Validate the company struct
	err := s.validator.Struct(company)
	if err != nil {
		return fmt.Errorf("%w: validation: %w", ErrCompanyService, err)
	}

	err = s.db.Create(company).Error
	if err != nil {
		return fmt.Errorf("%w: create: %w", ErrCompanyService, err)
	}

	return nil
}

func (s *Company) Get(id uuid.UUID) (*model.Company, error) {
	var company model.Company
	err := s.db.Where("id = ?", id).First(&company).Error
	if err != nil {
		return nil, fmt.Errorf("%w: get: %w", ErrCompanyService, err)
	}

	return &company, nil
}

func (s *Company) Update(company *model.Company) error {
	// Validate the company struct
	err := s.validator.Struct(company)
	if err != nil {
		return fmt.Errorf("%w: validation: %s", ErrCompanyService, err.Error())
	}

	err = s.db.Save(company).Error
	if err != nil {
		return fmt.Errorf("%w: update: %w", ErrCompanyService, err)
	}

	return nil
}

func (s *Company) Delete(id uuid.UUID) error {
	err := s.db.Where("id = ?", id).Delete(&model.Company{}).Error
	if err != nil {
		return fmt.Errorf("%w: delete: %w", ErrCompanyService, err)
	}

	return nil
}
