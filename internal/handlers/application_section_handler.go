package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

type ApplicationSectionHandler struct {
	cfg                    *config.Config
	applicationDetailsRepo *repository.ApplicationDetailsRepository
	applicationRepo        *repository.ApplicationRepository
}

func NewApplicationSectionHandler(cfg *config.Config) *ApplicationSectionHandler {
	return &ApplicationSectionHandler{
		cfg:                    cfg,
		applicationDetailsRepo: repository.NewApplicationDetailsRepository(),
		applicationRepo:        repository.NewApplicationRepository(),
	}
}

// SaveSection saves a specific section of the application
// @Summary Save application section
// @Description Save a specific section of the application (personal_info, address_info, education_history, family_info, financial_info, activities_skills)
// @Tags Application Sections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Application ID"
// @Param section_name path string true "Section name" Enums(personal_info, address_info, education_history, family_info, financial_info, activities_skills)
// @Param data body object true "Section data"
// @Success 200 {object} object{success=bool,message=string,data=object}
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Router /api/v1/applications/{id}/sections/{section_name} [post]
func (h *ApplicationSectionHandler) SaveSection(c *fiber.Ctx) error {
	// Get user_id from context
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

	sectionName := c.Params("section_name")
	if sectionName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Section name is required",
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

	// Verify application is still in draft status
	application, err := h.applicationRepo.GetByID(uint(applicationID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application",
		})
	}

	if application.ApplicationStatus != "draft" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot modify application that is not in draft status",
		})
	}

	// Route to appropriate save method based on section name
	switch sectionName {
	case "personal_info":
		return h.savePersonalInfo(c, applicationID)
	case "address_info":
		return h.saveAddressInfo(c, applicationID)
	case "education_history":
		return h.saveEducationHistory(c, applicationID)
	case "family_info":
		return h.saveFamilyInfo(c, applicationID)
	case "financial_info":
		return h.saveFinancialInfo(c, applicationID)
	case "activities_skills":
		return h.saveActivitiesSkills(c, applicationID)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid section name",
		})
	}
}

// PersonalInfoDTO is a data transfer object for personal info with simpler types
type PersonalInfoDTO struct {
	PrefixTH         *string `json:"prefix_th,omitempty"`
	PrefixEN         *string `json:"prefix_en,omitempty"`
	FirstNameTH      string  `json:"first_name_th"`
	LastNameTH       string  `json:"last_name_th"`
	FirstNameEN      *string `json:"first_name_en,omitempty"`
	LastNameEN       *string `json:"last_name_en,omitempty"`
	Email            string  `json:"email"`
	Phone            *string `json:"phone,omitempty"`
	LineID           *string `json:"line_id,omitempty"`
	CitizenID        *string `json:"citizen_id,omitempty"`
	StudentID        *string `json:"student_id,omitempty"`
	Faculty          *string `json:"faculty,omitempty"`
	Department       *string `json:"department,omitempty"`
	Major            *string `json:"major,omitempty"`
	YearLevel        *int    `json:"year_level,omitempty"`
	AdmissionType    *string `json:"admission_type,omitempty"`
	AdmissionDetails *string `json:"admission_details,omitempty"`
	DateOfBirth      *string `json:"date_of_birth,omitempty"`
	Gender           *string `json:"gender,omitempty"`
}

func (h *ApplicationSectionHandler) savePersonalInfo(c *fiber.Ctx, applicationID int) error {
	var dto PersonalInfoDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Convert DTO to model
	info := models.ApplicationPersonalInfo{
		ApplicationID: applicationID,
		FirstNameTH:   dto.FirstNameTH,
		LastNameTH:    dto.LastNameTH,
		Email:         dto.Email,
	}

	if dto.PrefixTH != nil {
		info.PrefixTH = sql.NullString{String: *dto.PrefixTH, Valid: true}
	}
	if dto.PrefixEN != nil {
		info.PrefixEN = sql.NullString{String: *dto.PrefixEN, Valid: true}
	}
	if dto.FirstNameEN != nil {
		info.FirstNameEN = sql.NullString{String: *dto.FirstNameEN, Valid: true}
	}
	if dto.LastNameEN != nil {
		info.LastNameEN = sql.NullString{String: *dto.LastNameEN, Valid: true}
	}
	if dto.Phone != nil {
		info.Phone = sql.NullString{String: *dto.Phone, Valid: true}
	}
	if dto.LineID != nil {
		info.LineID = sql.NullString{String: *dto.LineID, Valid: true}
	}
	if dto.CitizenID != nil {
		info.CitizenID = sql.NullString{String: *dto.CitizenID, Valid: true}
	}
	if dto.StudentID != nil {
		info.StudentID = sql.NullString{String: *dto.StudentID, Valid: true}
	}
	if dto.Faculty != nil {
		info.Faculty = sql.NullString{String: *dto.Faculty, Valid: true}
	}
	if dto.Department != nil {
		info.Department = sql.NullString{String: *dto.Department, Valid: true}
	}
	if dto.Major != nil {
		info.Major = sql.NullString{String: *dto.Major, Valid: true}
	}
	if dto.YearLevel != nil {
		info.YearLevel = sql.NullInt32{Int32: int32(*dto.YearLevel), Valid: true}
	}
	if dto.AdmissionType != nil {
		info.AdmissionType = sql.NullString{String: *dto.AdmissionType, Valid: true}
	}
	if dto.AdmissionDetails != nil {
		info.AdmissionDetails = sql.NullString{String: *dto.AdmissionDetails, Valid: true}
	}

	savedInfo, err := h.applicationDetailsRepo.SavePersonalInfo(&info)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save personal information",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Section saved successfully",
		"data": fiber.Map{
			"section":  "personal_info",
			"saved_at": savedInfo.UpdatedAt,
		},
	})
}

// AddressDTO is a data transfer object for address with simpler types
type AddressDTO struct {
	AddressType   string   `json:"address_type"`
	HouseNumber   *string  `json:"house_number,omitempty"`
	VillageNumber *string  `json:"village_number,omitempty"`
	Alley         *string  `json:"alley,omitempty"`
	Road          *string  `json:"road,omitempty"`
	Subdistrict   *string  `json:"subdistrict,omitempty"`
	District      *string  `json:"district,omitempty"`
	Province      *string  `json:"province,omitempty"`
	PostalCode    *string  `json:"postal_code,omitempty"`
	AddressLine1  *string  `json:"address_line1,omitempty"`
	AddressLine2  *string  `json:"address_line2,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	MapImageURL   *string  `json:"map_image_url,omitempty"`
}

func (h *ApplicationSectionHandler) saveAddressInfo(c *fiber.Ctx, applicationID int) error {
	var addressDTOs []AddressDTO
	if err := c.BodyParser(&addressDTOs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Convert DTOs to models
	addresses := make([]models.ApplicationAddress, len(addressDTOs))
	for i, dto := range addressDTOs {
		addresses[i] = models.ApplicationAddress{
			ApplicationID: applicationID,
			AddressType:   dto.AddressType,
		}

		if dto.HouseNumber != nil {
			addresses[i].HouseNumber = sql.NullString{String: *dto.HouseNumber, Valid: true}
		}
		if dto.VillageNumber != nil {
			addresses[i].VillageNumber = sql.NullString{String: *dto.VillageNumber, Valid: true}
		}
		if dto.Alley != nil {
			addresses[i].Alley = sql.NullString{String: *dto.Alley, Valid: true}
		}
		if dto.Road != nil {
			addresses[i].Road = sql.NullString{String: *dto.Road, Valid: true}
		}
		if dto.Subdistrict != nil {
			addresses[i].Subdistrict = sql.NullString{String: *dto.Subdistrict, Valid: true}
		}
		if dto.District != nil {
			addresses[i].District = sql.NullString{String: *dto.District, Valid: true}
		}
		if dto.Province != nil {
			addresses[i].Province = sql.NullString{String: *dto.Province, Valid: true}
		}
		if dto.PostalCode != nil {
			addresses[i].PostalCode = sql.NullString{String: *dto.PostalCode, Valid: true}
		}
		if dto.AddressLine1 != nil {
			addresses[i].AddressLine1 = sql.NullString{String: *dto.AddressLine1, Valid: true}
		}
		if dto.AddressLine2 != nil {
			addresses[i].AddressLine2 = sql.NullString{String: *dto.AddressLine2, Valid: true}
		}
		if dto.Latitude != nil {
			addresses[i].Latitude = sql.NullFloat64{Float64: *dto.Latitude, Valid: true}
		}
		if dto.Longitude != nil {
			addresses[i].Longitude = sql.NullFloat64{Float64: *dto.Longitude, Valid: true}
		}
		if dto.MapImageURL != nil {
			addresses[i].MapImageURL = sql.NullString{String: *dto.MapImageURL, Valid: true}
		}
	}

	savedAddresses, err := h.applicationDetailsRepo.SaveAddresses(uint(applicationID), addresses)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save address information",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Section saved successfully",
		"data": fiber.Map{
			"section":  "address_info",
			"saved_at": savedAddresses[0].UpdatedAt,
		},
	})
}

func (h *ApplicationSectionHandler) saveEducationHistory(c *fiber.Ctx, applicationID int) error {
	var education []models.ApplicationEducationHistory
	if err := c.BodyParser(&education); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	for i := range education {
		education[i].ApplicationID = applicationID
	}

	savedEducation, err := h.applicationDetailsRepo.SaveEducation(uint(applicationID), education)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save education history",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Section saved successfully",
		"data": fiber.Map{
			"section":  "education_history",
			"saved_at": savedEducation[0].UpdatedAt,
		},
	})
}

func (h *ApplicationSectionHandler) saveFamilyInfo(c *fiber.Ctx, applicationID int) error {
	var req struct {
		Members         []models.ApplicationFamilyMember   `json:"members"`
		Guardians       []models.ApplicationGuardian       `json:"guardians"`
		Siblings        []models.ApplicationSibling        `json:"siblings"`
		LivingSituation *models.ApplicationLivingSituation `json:"living_situation"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

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

	result, err := h.applicationDetailsRepo.SaveFamily(uint(applicationID), req.Members, req.Guardians, req.Siblings, req.LivingSituation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save family information",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Section saved successfully",
		"data": fiber.Map{
			"section":  "family_info",
			"saved_at": result["updated_at"],
		},
	})
}

func (h *ApplicationSectionHandler) saveFinancialInfo(c *fiber.Ctx, applicationID int) error {
	var req struct {
		FinancialInfo      *models.ApplicationFinancialInfo       `json:"financial_info"`
		Assets             []models.ApplicationAsset              `json:"assets"`
		ScholarshipHistory []models.ApplicationScholarshipHistory `json:"scholarship_history"`
		HealthInfo         *models.ApplicationHealthInfo          `json:"health_info"`
		FundingNeeds       *models.ApplicationFundingNeeds        `json:"funding_needs"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

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
		"message": "Section saved successfully",
		"data": fiber.Map{
			"section":  "financial_info",
			"saved_at": result["updated_at"],
		},
	})
}

func (h *ApplicationSectionHandler) saveActivitiesSkills(c *fiber.Ctx, applicationID int) error {
	var req struct {
		Activities []models.ApplicationActivity  `json:"activities"`
		References []models.ApplicationReference `json:"references"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	for i := range req.Activities {
		req.Activities[i].ApplicationID = applicationID
	}
	for i := range req.References {
		req.References[i].ApplicationID = applicationID
	}

	result, err := h.applicationDetailsRepo.SaveActivities(uint(applicationID), req.Activities, req.References)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save activities",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Section saved successfully",
		"data": fiber.Map{
			"section":  "activities_skills",
			"saved_at": result["updated_at"],
		},
	})
}

// verifyApplicationOwnership checks if the user owns the application
func (h *ApplicationSectionHandler) verifyApplicationOwnership(applicationID uint, userID uuid.UUID) error {
	application, err := h.applicationRepo.GetByID(applicationID)
	if err != nil {
		return err
	}

	// Get student ID from students table by user_id
	var studentID string
	err = database.DB.QueryRow(
		"SELECT student_id FROM students WHERE user_id = $1",
		userID,
	).Scan(&studentID)

	if err != nil {
		if err == sql.ErrNoRows {
			// If no student record exists, try matching with email
			userRepo := repository.NewUserRepository()
			user, userErr := userRepo.GetByID(userID)
			if userErr != nil {
				return userErr
			}

			if application.StudentID != user.Email {
				return fiber.NewError(fiber.StatusForbidden, "Unauthorized")
			}
			return nil
		}
		return err
	}

	if application.StudentID != studentID {
		return fiber.NewError(fiber.StatusForbidden, "Unauthorized")
	}

	return nil
}
