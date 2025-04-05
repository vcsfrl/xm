package example

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/vcsfrl/xm/internal/api/middleware"
	"github.com/vcsfrl/xm/internal/config"
	"github.com/vcsfrl/xm/internal/model"
	"io"
	"net/http"
)

func Run(config *config.Config, logger zerolog.Logger) {
	authResp, err := authenticate(middleware.LoginRequest{
		Username: config.AuthUser,
		Password: config.AuthPassword,
	}, config)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to authenticate")
		return
	}

	logger.Info().Msg("Authentication successful")
	logger.Info().Msgf("Token: %s", authResp.Token)

	// CREATE COMPANY
	//////////////////////////////////////////////////////////////////////////////////
	company := testCompany()
	companyJsonValue, err := json.Marshal(company)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal company")
		return
	}

	req, err := http.NewRequest("POST", baseUrl(config)+"/api/v1/company", bytes.NewBuffer(companyJsonValue))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authResp.Token))

	body, err := doRequest(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to do request")
		return
	}

	logger.Info().Msgf("Created Company: %s", string(body))

	// UPDATE COMPANY
	//////////////////////////////////////////////////////////////////////////////////
	var responseCompany model.Company
	err = json.Unmarshal(body, &responseCompany)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal company")
		return
	}

	responseCompany.Description = responseCompany.Description + " - updated"

	jsonValue, _ := json.Marshal(responseCompany)
	req, _ = http.NewRequest("PATCH", baseUrl(config)+"/api/v1/company/"+responseCompany.ID.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authResp.Token))
	body, err = doRequest(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to do request")
		return
	}

	logger.Info().Msgf("Updated Company: %s", string(body))

	// GET COMPANY
	//////////////////////////////////////////////////////////////////////////////////

	req, _ = http.NewRequest("GET", baseUrl(config)+"/api/v1/company/"+responseCompany.ID.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	body, err = doRequest(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to do request")
		return
	}

	logger.Info().Msgf("GET Company: %s", string(body))

	// DELETE COMPANY
	//////////////////////////////////////////////////////////////////////////////////
	req, _ = http.NewRequest("DELETE", baseUrl(config)+"/api/v1/company/"+responseCompany.ID.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authResp.Token))

	body, err = doRequest(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to do request")
		return
	}

	logger.Info().Msgf("Delete Response: %s", string(body))

}

func doRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func authenticate(login middleware.LoginRequest, config *config.Config) (*middleware.LoginResponse, error) {
	jsonValue, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseUrl(config)+"/api/v1/login", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	var loginResponse middleware.LoginResponse
	err = json.Unmarshal(body, &loginResponse)

	return &loginResponse, err
}

func testCompany() model.Company {
	return model.Company{
		Name:              "TestCompany - " + uuid.New().String(),
		Description:       "A test company",
		AmountOfEmployees: 10,
		Registered:        true,
		Type:              model.CompanyTypeCorporation,
	}
}

func baseUrl(config *config.Config) string {
	baseUrl := fmt.Sprintf("http://localhost:%s", config.AppPort)
	return baseUrl
}
