package handlers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type ApplicationDetailsHandler struct {
	cfg                     *config.Config
	applicationDetailsRepo  *repository.ApplicationDetailsRepository
	applicationRepo         *repository.ApplicationRepository
}

func NewApplicationDetailsHandler(cfg *config.Config) *ApplicationDetailsHandler {
	return &ApplicationDetailsHandler{
		cfg:                    cfg,
		applicationDetailsRepo: repository.NewApplicationDetailsRepository(),
		applicationRepo:        repository.NewApplicationRepository(),
	}
}

// SavePersonalInfo saves or updates personal information for an application
// @Summary Save personal information
// @Description Save or update personal information for an application (Student only)
// @Tags Application Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param info body models.ApplicationPersonalInfo true "Personal information"
// @Success 200 {object} object{success=bool,message=string,data=models.ApplicationPersonalInfo}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/personal-info [post]
func (h *ApplicationDetailsHandler) SavePersonalInfo(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only modify your own applications",
		})
	}

	var info models.ApplicationPersonalInfo
	if err := c.BodyParser(&info); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	info.ApplicationID = applicationID

	// Save or update personal info
	savedInfo, err := h.applicationDetailsRepo.SavePersonalInfo(&info)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save personal information",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "บันทึกข้อมูลส่วนตัวเรียบร้อยแล้ว",
		"data":    savedInfo,
	})
}

// SaveAddresses saves or updates addresses for an application
// @Summary Save addresses
// @Description Save or update addresses for an application (Student only)
// @Tags Application Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param addresses body []models.ApplicationAddress true "Addresses"
// @Success 200 {object} object{success=bool,message=string,data=[]models.ApplicationAddress}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/addresses [post]
func (h *ApplicationDetailsHandler) SaveAddresses(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only modify your own applications",
		})
	}

	var addresses []models.ApplicationAddress
	if err := c.BodyParser(&addresses); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Set application ID for all addresses
	for i := range addresses {
		addresses[i].ApplicationID = applicationID
	}

	// Save addresses
	savedAddresses, err := h.applicationDetailsRepo.SaveAddresses(uint(applicationID), addresses)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save addresses",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "บันทึกข้อมูลที่อยู่เรียบร้อยแล้ว",
		"data":    savedAddresses,
	})
}

// SaveEducation saves or updates education history for an application
// @Summary Save education history
// @Description Save or update education history for an application (Student only)
// @Tags Application Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param education body []models.ApplicationEducationHistory true "Education history"
// @Success 200 {object} object{success=bool,message=string,data=[]models.ApplicationEducationHistory}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/education [post]
func (h *ApplicationDetailsHandler) SaveEducation(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only modify your own applications",
		})
	}

	var education []models.ApplicationEducationHistory
	if err := c.BodyParser(&education); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Set application ID for all records
	for i := range education {
		education[i].ApplicationID = applicationID
	}

	// Save education history
	savedEducation, err := h.applicationDetailsRepo.SaveEducation(uint(applicationID), education)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save education history",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "บันทึกข้อมูลการศึกษาเรียบร้อยแล้ว",
		"data":    savedEducation,
	})
}

// SaveFamily saves or updates family information for an application
// @Summary Save family information
// @Description Save or update family information for an application (Student only)
// @Tags Application Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param family body object{members=[]models.ApplicationFamilyMember,guardians=[]models.ApplicationGuardian,siblings=[]models.ApplicationSibling,living_situation=models.ApplicationLivingSituation} true "Family information"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/family [post]
func (h *ApplicationDetailsHandler) SaveFamily(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only modify your own applications",
		})
	}

	var req struct {
		Members         []models.ApplicationFamilyMember    `json:"members"`
		Guardians       []models.ApplicationGuardian        `json:"guardians"`
		Siblings        []models.ApplicationSibling         `json:"siblings"`
		LivingSituation *models.ApplicationLivingSituation  `json:"living_situation"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Set application ID for all records
	for i := range req.Members {
		req.Members[i].ApplicationID = applicationID
	}
	for i := range req.Guardians {
		req.Guardians[i].ApplicationID = applicationID
	}
	for i := range req.Siblings {
		req.Siblings[i].ApplicationID = applicationID
	}
	if req.LivingSituation != nil {
		req.LivingSituation.ApplicationID = applicationID
	}

	// Save family information
	result, err := h.applicationDetailsRepo.SaveFamily(uint(applicationID), req.Members, req.Guardians, req.Siblings, req.LivingSituation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save family information",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "บันทึกข้อมูลครอบครัวเรียบร้อยแล้ว",
		"data":    result,
	})
}

// SaveFinancial saves or updates financial information for an application
// @Summary Save financial information
// @Description Save or update financial information for an application (Student only)
// @Tags Application Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param financial body object{financial_info=models.ApplicationFinancialInfo,assets=[]models.ApplicationAsset,scholarship_history=[]models.ApplicationScholarshipHistory,health_info=models.ApplicationHealthInfo,funding_needs=models.ApplicationFundingNeeds} true "Financial information"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/financial [post]
func (h *ApplicationDetailsHandler) SaveFinancial(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only modify your own applications",
		})
	}

	var req struct {
		FinancialInfo      *models.ApplicationFinancialInfo      `json:"financial_info"`
		Assets             []models.ApplicationAsset             `json:"assets"`
		ScholarshipHistory []models.ApplicationScholarshipHistory `json:"scholarship_history"`
		HealthInfo         *models.ApplicationHealthInfo         `json:"health_info"`
		FundingNeeds       *models.ApplicationFundingNeeds       `json:"funding_needs"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Set application ID for all records
	if req.FinancialInfo != nil {
		req.FinancialInfo.ApplicationID = applicationID
	}
	for i := range req.Assets {
		req.Assets[i].ApplicationID = applicationID
	}
	for i := range req.ScholarshipHistory {
		req.ScholarshipHistory[i].ApplicationID = applicationID
	}
	if req.HealthInfo != nil {
		req.HealthInfo.ApplicationID = applicationID
	}
	if req.FundingNeeds != nil {
		req.FundingNeeds.ApplicationID = applicationID
	}

	// Save financial information
	result, err := h.applicationDetailsRepo.SaveFinancial(
		uint(applicationID),
		req.FinancialInfo,
		req.Assets,
		req.ScholarshipHistory,
		req.HealthInfo,
		req.FundingNeeds,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save financial information",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "บันทึกข้อมูลการเงินเรียบร้อยแล้ว",
		"data":    result,
	})
}

// SaveActivities saves or updates activities and references for an application
// @Summary Save activities
// @Description Save or update activities and references for an application (Student only)
// @Tags Application Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param activities body object{activities=[]models.ApplicationActivity,references=[]models.ApplicationReference} true "Activities and references"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/activities [post]
func (h *ApplicationDetailsHandler) SaveActivities(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only modify your own applications",
		})
	}

	var req struct {
		Activities []models.ApplicationActivity  `json:"activities"`
		References []models.ApplicationReference `json:"references"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Set application ID for all records
	for i := range req.Activities {
		req.Activities[i].ApplicationID = applicationID
	}
	for i := range req.References {
		req.References[i].ApplicationID = applicationID
	}

	// Save activities
	result, err := h.applicationDetailsRepo.SaveActivities(uint(applicationID), req.Activities, req.References)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save activities",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "บันทึกข้อมูลกิจกรรมเรียบร้อยแล้ว",
		"data":    result,
	})
}

// SaveCompleteForm saves all application details at once
// @Summary Save complete application form
// @Description Save all application details at once (Student only)
// @Tags Application Details
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param form body models.CompleteApplicationForm true "Complete application form"
// @Success 200 {object} object{success=bool,message=string,data=models.CompleteApplicationForm}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/complete-form [post]
func (h *ApplicationDetailsHandler) SaveCompleteForm(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only modify your own applications",
		})
	}

	var form models.CompleteApplicationForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Set application ID for all records
	if form.PersonalInfo != nil {
		form.PersonalInfo.ApplicationID = applicationID
	}
	for i := range form.Addresses {
		form.Addresses[i].ApplicationID = applicationID
	}
	for i := range form.EducationHistory {
		form.EducationHistory[i].ApplicationID = applicationID
	}
	for i := range form.FamilyMembers {
		form.FamilyMembers[i].ApplicationID = applicationID
	}
	for i := range form.Assets {
		form.Assets[i].ApplicationID = applicationID
	}
	for i := range form.Guardians {
		form.Guardians[i].ApplicationID = applicationID
	}
	for i := range form.Siblings {
		form.Siblings[i].ApplicationID = applicationID
	}
	if form.LivingSituation != nil {
		form.LivingSituation.ApplicationID = applicationID
	}
	if form.FinancialInfo != nil {
		form.FinancialInfo.ApplicationID = applicationID
	}
	for i := range form.ScholarshipHistory {
		form.ScholarshipHistory[i].ApplicationID = applicationID
	}
	for i := range form.Activities {
		form.Activities[i].ApplicationID = applicationID
	}
	for i := range form.References {
		form.References[i].ApplicationID = applicationID
	}
	if form.HealthInfo != nil {
		form.HealthInfo.ApplicationID = applicationID
	}
	if form.FundingNeeds != nil {
		form.FundingNeeds.ApplicationID = applicationID
	}
	for i := range form.HouseDocuments {
		form.HouseDocuments[i].ApplicationID = applicationID
	}
	for i := range form.IncomeCertificates {
		form.IncomeCertificates[i].ApplicationID = applicationID
	}

	// Save complete form
	savedForm, err := h.applicationDetailsRepo.SaveCompleteForm(uint(applicationID), &form)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save application form",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "บันทึกข้อมูลใบสมัครทั้งหมดเรียบร้อยแล้ว",
		"data":    savedForm,
	})
}

// GetCompleteForm retrieves all application details
// @Summary Get complete application form
// @Description Get all application details (Student/Admin/Officer)
// @Tags Application Details
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} object{success=bool,data=models.CompleteApplicationForm}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/complete-form [get]
func (h *ApplicationDetailsHandler) GetCompleteForm(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Check if user has access to this application
	roles := c.Locals("roles").([]string)
	isAdmin := false
	for _, role := range roles {
		if role == "admin" || role == "officer" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		// Student can only view their own application
		if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Application not found",
				})
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You can only view your own applications",
			})
		}
	}

	// Get complete form
	form, err := h.applicationDetailsRepo.GetCompleteForm(uint(applicationID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve application details",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    form,
	})
}

// SubmitApplication submits the application for review
// @Summary Submit application
// @Description Submit the application for review (Student only)
// @Tags Application Details
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Success 200 {object} object{success=bool,message=string}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/submit [put]
func (h *ApplicationDetailsHandler) SubmitApplication(c *fiber.Ctx) error {
	// Get user_id from context - handle both UUID and string formats
	var userID uuid.UUID
	userIDValue := c.Locals("user_id")

	switch v := userIDValue.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}
		userID = parsed
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	applicationIDStr := c.Params("id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid application ID",
		})
	}

	// Verify application ownership
	if err := h.verifyApplicationOwnership(uint(applicationID), userID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only submit your own applications",
		})
	}

	// Get application
	application, err := h.applicationRepo.GetByID(uint(applicationID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Application not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application",
		})
	}

	// Check if application is in draft status
	if application.ApplicationStatus != "draft" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Application has already been submitted",
		})
	}

	// Validate that required information is complete
	form, err := h.applicationDetailsRepo.GetCompleteForm(uint(applicationID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Please complete all required information before submitting",
		})
	}

	// Validate required fields
	if form.PersonalInfo == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Personal information is required",
		})
	}

	// Submit application
	if err := h.applicationRepo.Submit(uint(applicationID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to submit application",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ส่งใบสมัครเรียบร้อยแล้ว รอการตรวจสอบจากเจ้าหน้าที่",
	})
}

// verifyApplicationOwnership checks if the user owns the application
func (h *ApplicationDetailsHandler) verifyApplicationOwnership(applicationID uint, userID uuid.UUID) error {
	fmt.Printf("Debug - verifyApplicationOwnership called: appID=%d, userID=%s\n", applicationID, userID.String())

	application, err := h.applicationRepo.GetByID(applicationID)
	if err != nil {
		fmt.Printf("Debug - Application not found: %v\n", err)
		return err
	}

	fmt.Printf("Debug - Application found: ID=%d, StudentID=%s\n", application.ApplicationID, application.StudentID)

	// Get student ID from students table by user_id
	var studentID string
	err = database.DB.QueryRow(
		"SELECT student_id FROM students WHERE user_id = $1",
		userID,
	).Scan(&studentID)

	if err != nil {
		if err == sql.ErrNoRows {
			// If no student record exists, try matching with email for backwards compatibility
			fmt.Printf("Debug - No student record found for userID=%s, trying email fallback\n", userID.String())
			userRepo := repository.NewUserRepository()
			user, userErr := userRepo.GetByID(userID)
			if userErr != nil {
				fmt.Printf("Debug - User not found: %v\n", userErr)
				return userErr
			}

			fmt.Printf("Debug - User found: ID=%s, Email=%s\n", user.UserID.String(), user.Email)

			// Try email comparison as fallback
			if application.StudentID != user.Email {
				fmt.Printf("Debug - Ownership verification FAILED (email fallback): app.StudentID(%s) != user.Email(%s)\n", application.StudentID, user.Email)
				return fiber.NewError(fiber.StatusForbidden, "Unauthorized")
			}

			fmt.Printf("Debug - Ownership verification SUCCESS (email fallback)\n")
			return nil
		}
		fmt.Printf("Debug - Database error: %v\n", err)
		return err
	}

	fmt.Printf("Debug - Student found: StudentID=%s, UserID=%s\n", studentID, userID.String())

	// Compare student IDs
	if application.StudentID != studentID {
		fmt.Printf("Debug - Ownership verification FAILED: app.StudentID(%s) != studentID(%s)\n", application.StudentID, studentID)
		return fiber.NewError(fiber.StatusForbidden, "Unauthorized")
	}

	fmt.Printf("Debug - Ownership verification SUCCESS\n")
	return nil
}
