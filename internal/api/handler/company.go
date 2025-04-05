package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vcsfrl/xm/internal/model"
	"github.com/vcsfrl/xm/internal/service"
	"net/http"
)

type CompanyHandler struct {
	company *service.Company
}

func NewCompanyHandler(company *service.Company) *CompanyHandler {
	return &CompanyHandler{company: company}
}

func (ch *CompanyHandler) Create(c *gin.Context) {
	var company model.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ch.company.Create(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, company)
}

func (ch *CompanyHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var company *model.Company
	var err error

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	company, err = ch.company.Get(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting company"})
		return
	}

	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (ch *CompanyHandler) Update(c *gin.Context) {
	var company *model.Company
	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	company, err = ch.company.Get(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting company"})
		return
	}

	fmt.Println(company)

	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	if err := c.ShouldBindJSON(company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	company.ID = uuid

	if err := ch.company.Update(company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating company"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (ch *CompanyHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := ch.company.Delete(uuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted"})
}
