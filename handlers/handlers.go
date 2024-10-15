package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imnotjin/fetch-backend-challenge/models"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	DB *gorm.DB
}

type AddPointsRequest struct {
	Payer     string    `json:"payer" binding:"required"`
	Points    int       `json:"points" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

type SpendPointsRequest struct {
	Points int `json:"points" binding:"required"`
}

type SpendPointsResponse struct {
	Payer  string `json:"payer"`
	Points int    `json:"points"`
}

// AddPoints godoc
// @Summary Add points to a user's account
// @Description Add points to a user's account for a specific payer
// @Tags points
// @Accept json
// @Produce json
// @Param request body AddPointsRequest true "Add points request"
// @Success 200 {string} string "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /add [post]
func (h *Handler) AddPoints(c *gin.Context) {
	var req AddPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := models.Transaction{
		Payer:     req.Payer,
		Points:    req.Points,
		Timestamp: req.Timestamp,
	}

	if err := h.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add points"})
		return
	}

	c.Status(http.StatusOK)
}

// SpendPoints godoc
// @Summary Spend points from user's account
// @Description Spend points from user's account following the given rules
// @Tags points
// @Accept json
// @Produce json
// @Param request body SpendPointsRequest true "Spend points request"
// @Success 200 {array} SpendPointsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /spend [post]
func (h *Handler) SpendPoints(c *gin.Context) {
	var req SpendPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transactions []models.Transaction
	if err := h.DB.Order("timestamp").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	totalPoints := 0
	for _, t := range transactions {
		totalPoints += t.Points
	}

	if totalPoints < req.Points {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough points"})
		return
	}

	pointsToSpend := req.Points
	spentPoints := make(map[string]int)

	for i := 0; i < len(transactions) && pointsToSpend > 0; i++ {
		t := &transactions[i]
		if t.Points <= pointsToSpend {
			spentPoints[t.Payer] += t.Points
			pointsToSpend -= t.Points
			t.Points = 0
		} else {
			spentPoints[t.Payer] += pointsToSpend
			t.Points -= pointsToSpend
			pointsToSpend = 0
		}
	}

	var response []SpendPointsResponse
	for payer, points := range spentPoints {
		response = append(response, SpendPointsResponse{Payer: payer, Points: -points})
	}

	// Update the database
	if err := h.DB.Save(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update points"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetBalance godoc
// @Summary Get user's point balance
// @Description Get the current point balance for each payer
// @Tags points
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} ErrorResponse
// @Router /balance [get]
func (h *Handler) GetBalance(c *gin.Context) {
	var transactions []models.Transaction
	if err := h.DB.Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	balance := make(map[string]int)
	for _, t := range transactions {
		balance[t.Payer] += t.Points
	}

	c.JSON(http.StatusOK, balance)
}
