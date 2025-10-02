package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
	"scholarship-system/internal/repository"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	repo *repository.PaymentRepository
	cfg  *config.Config
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(cfg *config.Config) *PaymentHandler {
	return &PaymentHandler{
		repo: repository.NewPaymentRepository(database.DB),
		cfg:  cfg,
	}
}

// CreateTransaction creates a new payment transaction
// @Summary Create payment transaction
// @Description Create a new payment transaction for scholarship allocation
// @Tags payments
// @Accept json
// @Produce json
// @Param transaction body models.PaymentTransaction true "Transaction data"
// @Success 201 {object} models.PaymentTransaction
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/payments/transactions [post]
func (h *PaymentHandler) CreateTransaction(c *fiber.Ctx) error {
	var tx models.PaymentTransaction
	if err := c.BodyParser(&tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Generate UUID if not provided
	if tx.TransactionID == uuid.Nil {
		tx.TransactionID = uuid.New()
	}

	// Set defaults
	if tx.PaymentStatus == "" {
		tx.PaymentStatus = "pending"
	}

	if err := h.repo.CreateTransaction(&tx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tx)
}

// GetTransaction retrieves a transaction by ID
// @Summary Get payment transaction
// @Description Get payment transaction details by ID
// @Tags payments
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} models.PaymentTransaction
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/payments/transactions/{id} [get]
func (h *PaymentHandler) GetTransaction(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	tx, err := h.repo.GetTransactionByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	return c.JSON(tx)
}

// GetAllocationTransactions gets all transactions for an allocation
// @Summary Get allocation transactions
// @Description Get all payment transactions for a scholarship allocation
// @Tags payments
// @Produce json
// @Param allocation_id path int true "Allocation ID"
// @Success 200 {array} models.PaymentTransaction
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/payments/allocations/{allocation_id}/transactions [get]
func (h *PaymentHandler) GetAllocationTransactions(c *fiber.Ctx) error {
	allocationID, err := c.ParamsInt("allocation_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid allocation ID",
		})
	}

	transactions, err := h.repo.GetTransactionsByAllocation(allocationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	return c.JSON(transactions)
}

// UpdateTransactionStatus updates transaction status
// @Summary Update transaction status
// @Description Update the status of a payment transaction
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param status body map[string]string true "Status update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/payments/transactions/{id}/status [put]
func (h *PaymentHandler) UpdateTransactionStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID",
		})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.repo.UpdateTransactionStatus(id, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update status",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Status updated successfully",
	})
}

// CreateDisbursementSchedule creates a disbursement schedule
// @Summary Create disbursement schedule
// @Description Create a disbursement schedule for installment payments
// @Tags payments
// @Accept json
// @Produce json
// @Param schedule body models.DisbursementSchedule true "Schedule data"
// @Success 201 {object} models.DisbursementSchedule
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/payments/disbursements [post]
func (h *PaymentHandler) CreateDisbursementSchedule(c *fiber.Ctx) error {
	var schedule models.DisbursementSchedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if schedule.ScheduleID == uuid.Nil {
		schedule.ScheduleID = uuid.New()
	}

	if schedule.Status == "" {
		schedule.Status = "scheduled"
	}

	if err := h.repo.CreateDisbursementSchedule(&schedule); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create schedule",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(schedule)
}

// GetDisbursementSchedules gets disbursement schedules
// @Summary Get disbursement schedules
// @Description Get all disbursement schedules for an allocation
// @Tags payments
// @Produce json
// @Param allocation_id path int true "Allocation ID"
// @Success 200 {array} models.DisbursementSchedule
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/payments/allocations/{allocation_id}/disbursements [get]
func (h *PaymentHandler) GetDisbursementSchedules(c *fiber.Ctx) error {
	allocationID, err := c.ParamsInt("allocation_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid allocation ID",
		})
	}

	schedules, err := h.repo.GetDisbursementSchedules(allocationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve schedules",
		})
	}

	return c.JSON(schedules)
}

// GetPendingDisbursements gets pending disbursements
// @Summary Get pending disbursements
// @Description Get all disbursements that are due for payment
// @Tags payments
// @Produce json
// @Success 200 {array} models.DisbursementSchedule
// @Security BearerAuth
// @Router /api/v1/payments/disbursements/pending [get]
func (h *PaymentHandler) GetPendingDisbursements(c *fiber.Ctx) error {
	schedules, err := h.repo.GetPendingDisbursements(time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve pending disbursements",
		})
	}

	return c.JSON(schedules)
}

// GetPaymentMethods gets all payment methods
// @Summary Get payment methods
// @Description Get all available payment methods
// @Tags payments
// @Produce json
// @Success 200 {array} models.PaymentMethod
// @Security BearerAuth
// @Router /api/v1/payments/methods [get]
func (h *PaymentHandler) GetPaymentMethods(c *fiber.Ctx) error {
	methods, err := h.repo.GetPaymentMethods()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve payment methods",
		})
	}

	return c.JSON(methods)
}
