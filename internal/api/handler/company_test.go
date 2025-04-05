package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/xm/internal/api/validator"
	db2 "github.com/vcsfrl/xm/internal/db"
	"github.com/vcsfrl/xm/internal/model"
	"github.com/vcsfrl/xm/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompanyHandler(t *testing.T) {
	suite.Run(t, new(CompanyHandlerSuite))
}

type CompanyHandlerSuite struct {
	suite.Suite
	companyService *service.Company
	handler        *CompanyHandler
	router         *gin.Engine
	logger         zerolog.Logger
}

func (suite *CompanyHandlerSuite) SetupTest() {
	suite.logger = zerolog.Nop()
	db, err := db2.InitTestSqlite(suite.logger)
	suite.NoError(err)
	suite.NotNil(db)

	suite.companyService = service.NewCompanyService(db, validator.CompanyValidator(suite.logger))
	suite.handler = NewCompanyHandler(suite.companyService)

	suite.router = gin.Default()
	suite.router.POST("/company", suite.handler.Create)
	suite.router.GET("/company/:id", suite.handler.Get)
	suite.router.PATCH("/company/:id", suite.handler.Update)
	suite.router.DELETE("/company/:id", suite.handler.Delete)
}

func (suite *CompanyHandlerSuite) TestCreateCompany() {
	company := model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              model.CompanyTypeCorporation,
	}
	jsonValue, _ := json.Marshal(company)
	req, _ := http.NewRequest("POST", "/company", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var responseCompany model.Company
	err := json.Unmarshal(w.Body.Bytes(), &responseCompany)
	suite.NoError(err)

	suite.Equal(company.Name, responseCompany.Name)
}

func (suite *CompanyHandlerSuite) TestGetCompany() {
	company := model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              model.CompanyTypeCorporation,
	}
	err := suite.companyService.Create(&company)
	suite.NoError(err)

	req, _ := http.NewRequest("GET", "/company/"+company.ID.String(), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var responseCompany model.Company
	err = json.Unmarshal(w.Body.Bytes(), &responseCompany)
	suite.NoError(err)

	suite.Equal(company.Name, responseCompany.Name)
}

func (suite *CompanyHandlerSuite) TestUpdateCompany() {
	company := model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              model.CompanyTypeCorporation,
	}
	err := suite.companyService.Create(&company)
	suite.NoError(err)

	updatedCompany := model.Company{
		Name:              "UpdatedCompany",
		Description:       "An updated test company",
		AmountOfEmployees: 20,
		Registered:        false,
		Type:              model.CompanyTypeNonProfit,
	}
	jsonValue, _ := json.Marshal(updatedCompany)
	req, _ := http.NewRequest("PATCH", "/company/"+company.ID.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var responseCompany model.Company
	err = json.Unmarshal(w.Body.Bytes(), &responseCompany)
	suite.NoError(err)

	suite.Equal(updatedCompany.Name, responseCompany.Name)
	suite.Equal(updatedCompany.Description, responseCompany.Description)
	suite.Equal(updatedCompany.AmountOfEmployees, responseCompany.AmountOfEmployees)
	suite.Equal(updatedCompany.Registered, responseCompany.Registered)
	suite.Equal(updatedCompany.Type, responseCompany.Type)
	suite.Equal(company.ID.String(), responseCompany.ID.String())
}

func (suite *CompanyHandlerSuite) TestUpdatePartialCompany() {
	company := model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              model.CompanyTypeCorporation,
	}
	err := suite.companyService.Create(&company)
	suite.NoError(err)

	updatedCompany := model.Company{
		Name: "UpdatedCompany",
	}
	jsonValue, _ := json.Marshal(updatedCompany)
	req, _ := http.NewRequest("PATCH", "/company/"+company.ID.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var responseCompany model.Company
	err = json.Unmarshal(w.Body.Bytes(), &responseCompany)
	suite.NoError(err)

	suite.Equal(updatedCompany.Name, responseCompany.Name)
	suite.Equal(company.Description, responseCompany.Description)
}

func (suite *CompanyHandlerSuite) TestDeleteCompany() {
	company := model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              model.CompanyTypeCorporation,
	}
	err := suite.companyService.Create(&company)
	suite.NoError(err)

	req, _ := http.NewRequest("DELETE", "/company/"+company.ID.String(), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal("Company deleted", response["message"])
}
