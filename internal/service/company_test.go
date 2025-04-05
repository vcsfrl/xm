package service

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/xm/internal/db"
	"github.com/vcsfrl/xm/internal/model"
	"github.com/vcsfrl/xm/internal/validator"
	"gorm.io/gorm"
	"testing"
)

func TestCompany(t *testing.T) {
	suite.Run(t, new(CompanyFixture))
}

type CompanyFixture struct {
	suite.Suite

	db     *gorm.DB
	logger zerolog.Logger
}

func (cf *CompanyFixture) SetupTest() {
	var err error
	cf.db, err = db.InitTestSqlite(zerolog.Nop())
	cf.NoError(err)
	cf.logger = zerolog.Nop()
}

func (cf *CompanyFixture) TestCreate() {
	service := NewCompanyService(cf.db, validator.CompanyValidator(cf.logger))

	company := &model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              "Corporations",
	}

	err := service.Create(company)
	cf.NoError(err)
	cf.NotNil(company.ID)
}

func (cf *CompanyFixture) TestGet() {
	service := NewCompanyService(cf.db, validator.CompanyValidator(cf.logger))

	company := &model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              "Corporations",
	}
	cf.db.Create(company)

	result, err := service.Get(company.ID)
	cf.NoError(err)
	cf.Equal(company.Name, result.Name)
}

func (cf *CompanyFixture) TestUpdate() {
	service := NewCompanyService(cf.db, validator.CompanyValidator(cf.logger))

	company := &model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              "Corporations",
	}
	cf.db.Create(company)

	company.Name = "UpdatedCompany"
	company.Description = "An updated test company"
	company.AmountOfEmployees = 20
	company.Registered = false
	company.Type = model.CompanyTypeSoleProprietorship
	err := service.Update(company)
	cf.NoError(err)

	result, _ := service.Get(company.ID)
	cf.Equal("UpdatedCompany", result.Name)
	cf.Equal("An updated test company", result.Description)
	cf.Equal(20, result.AmountOfEmployees)
	cf.Equal(false, result.Registered)
	cf.Equal(model.CompanyTypeSoleProprietorship, result.Type)
}

func (cf *CompanyFixture) TestDelete() {
	service := NewCompanyService(cf.db, validator.CompanyValidator(cf.logger))

	company := &model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              "Corporations",
	}
	cf.db.Create(company)

	err := service.Delete(company.ID)
	cf.NoError(err)

	result, err := service.Get(company.ID)
	cf.Error(err)
	cf.Nil(result)
}
