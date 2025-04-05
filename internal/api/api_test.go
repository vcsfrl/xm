package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/xm/internal/api/middleware"
	"github.com/vcsfrl/xm/internal/config"
	db2 "github.com/vcsfrl/xm/internal/db"
	"github.com/vcsfrl/xm/internal/model"
	"github.com/vcsfrl/xm/internal/service"
	"github.com/vcsfrl/xm/internal/validator"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestAPiSuite(t *testing.T) {
	suite.Run(t, new(RestApiTestSuite))
}

type RestApiTestSuite struct {
	suite.Suite
	logger         zerolog.Logger
	companyService *service.Company
	ctx            context.Context
	companyApi     *RestApi
	config         *config.Config
}

func (suite *RestApiTestSuite) SetupTest() {
	suite.logger = zerolog.Nop()
	db, err := db2.InitTestSqlite()
	suite.NoError(err)
	suite.NotNil(db)

	suite.ctx = context.Background()
	suite.config = &config.Config{
		AppPort: "1234", AuthJwtSecret: "secret", AuthUser: "admin", AuthPassword: "admin",
		RateLimit: 1000, RateBurst: 100}
	suite.companyApi = NewRestApi(suite.ctx, suite.logger, suite.config, db)
	suite.companyService = service.NewCompanyService(db, validator.CompanyValidator(suite.logger))
}

func (suite *RestApiTestSuite) TearDownTest() {
	err := suite.companyApi.Close()
	suite.NoError(err)
}

func (suite *RestApiTestSuite) TestApi_Login() {
	loginResponse := suite.authenticate(suite.loginRequest())

	suite.Equal(http.StatusOK, loginResponse.Code)
	suite.NotEmpty(loginResponse.ExpiresAt)
	suite.NotEmpty(loginResponse.Token)
}

func (suite *RestApiTestSuite) TestCreateCompany_Unauthorized() {
	jsonValue, err := json.Marshal(suite.testCompany())
	suite.NoError(err)

	// Test unauthorized request
	req, err := http.NewRequest("POST", "/api/v1/company", bytes.NewBuffer(jsonValue))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	router, err := suite.companyApi.BuildRouter()
	suite.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// If user is not logged in check that we get response code 401
	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *RestApiTestSuite) TestCreateCompany_Authorized() {
	loginResponse := suite.authenticate(suite.loginRequest())

	company := suite.testCompany()
	companyJsonValue, err := json.Marshal(company)
	suite.NoError(err)

	req, err := http.NewRequest("POST", "/api/v1/company", bytes.NewBuffer(companyJsonValue))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginResponse.Token))

	w := httptest.NewRecorder()
	router, err := suite.companyApi.BuildRouter()
	suite.NoError(err)
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var responseCompany model.Company
	err = json.Unmarshal(w.Body.Bytes(), &responseCompany)
	suite.NoError(err)

	suite.Equal(company.Name, responseCompany.Name)
}

func (suite *RestApiTestSuite) TestGetCompany() {
	company := suite.testCompany()
	err := suite.companyService.Create(&company)
	suite.NoError(err)

	req, _ := http.NewRequest("GET", "/api/v1/company/"+company.ID.String(), nil)
	w := httptest.NewRecorder()
	router, err := suite.companyApi.BuildRouter()
	suite.NoError(err)
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var responseCompany model.Company
	err = json.Unmarshal(w.Body.Bytes(), &responseCompany)
	suite.NoError(err)

	suite.Equal(company.Name, responseCompany.Name)
	suite.Equal(company.Description, responseCompany.Description)
	suite.Equal(company.AmountOfEmployees, responseCompany.AmountOfEmployees)
	suite.Equal(company.Registered, responseCompany.Registered)
	suite.Equal(company.Type, responseCompany.Type)
	suite.Equal(company.ID, responseCompany.ID)
}

func (suite *RestApiTestSuite) TestUpdateCompany_Unauthorized() {
	company := suite.testCompany()
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
	req, _ := http.NewRequest("PATCH", "/api/v1/company/"+company.ID.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router, err := suite.companyApi.BuildRouter()
	suite.NoError(err)
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *RestApiTestSuite) TestUpdateCompany_Authorized() {
	loginResponse := suite.authenticate(suite.loginRequest())
	company := suite.testCompany()
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
	req, _ := http.NewRequest("PATCH", "/api/v1/company/"+company.ID.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginResponse.Token))

	w := httptest.NewRecorder()
	router, err := suite.companyApi.BuildRouter()
	suite.NoError(err)
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var responseCompany model.Company
	err = json.Unmarshal(w.Body.Bytes(), &responseCompany)
	suite.NoError(err)

	suite.Equal(updatedCompany.Name, responseCompany.Name)
	suite.Equal(updatedCompany.Description, responseCompany.Description)
	suite.Equal(updatedCompany.AmountOfEmployees, responseCompany.AmountOfEmployees)
	suite.Equal(updatedCompany.Registered, responseCompany.Registered)
	suite.Equal(updatedCompany.Type, responseCompany.Type)
}

func (suite *RestApiTestSuite) TestDeleteCompany_Authorized() {
	loginResponse := suite.authenticate(suite.loginRequest())
	company := suite.testCompany()
	err := suite.companyService.Create(&company)
	suite.NoError(err)

	req, _ := http.NewRequest("DELETE", "/api/v1/company/"+company.ID.String(), nil)
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginResponse.Token))
	router, err := suite.companyApi.BuildRouter()
	suite.NoError(err)
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal("Company deleted", response["message"])
}

func (suite *RestApiTestSuite) TestDeleteCompany_Unauthorized() {
	company := suite.testCompany()
	err := suite.companyService.Create(&company)
	suite.NoError(err)

	req, _ := http.NewRequest("DELETE", "/api/v1/company/"+company.ID.String(), nil)
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")
	router, err := suite.companyApi.BuildRouter()
	suite.NoError(err)
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *RestApiTestSuite) testCompany() model.Company {
	return model.Company{
		Name:              "TestCompany",
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              model.CompanyTypeCorporation,
	}
}

func (suite *RestApiTestSuite) loginRequest() middleware.LoginRequest {
	return middleware.LoginRequest{
		Username: suite.config.AuthUser,
		Password: suite.config.AuthPassword,
	}
}

func (suite *RestApiTestSuite) authenticate(login middleware.LoginRequest) middleware.LoginResponse {
	jsonValue, err := json.Marshal(login)
	suite.NoError(err)
	req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonValue))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router, err := suite.companyApi.BuildRouter()
	suite.NoError(err)
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	var loginResponse middleware.LoginResponse
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	suite.NoError(err)
	return loginResponse
}
