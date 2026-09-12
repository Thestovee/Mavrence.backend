package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbPool struct {
	db *pgxpool.Pool
}

func NewDbPool(db *pgxpool.Pool) *DbPool {
	return &DbPool{
		db: db,
	}
}

type HealthData struct {
	ID            int       `json:"id"`
	Date          time.Time `json:"date"`
	MentalState   int       `json:"mental_state"`
	PhysicalState int       `json:"physical_state"`
	Note          *string   `json:"note,omitempty"`
}

type HealthDataRequest struct {
	MentalState   int `json:"mental_state" binding:"required"`
	PhysicalState int `json:"physical_state" binding:"required"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (h *DbPool) GetRecent(c *gin.Context) {
	ctx := c.Request.Context()
	rows, err := h.db.Query(ctx, "SELECT id, date, mentalState, physicalState FROM MavHealth ORDER BY date DESC LIMIT 7")
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to query health data from database ",
			Error:   err.Error()})
	}
	defer rows.Close()

	var records []HealthData

	for rows.Next() {
		var h HealthData

		err := rows.Scan(&h.ID, &h.Date, &h.MentalState, &h.PhysicalState)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to scan rows",
				Error:   err.Error(),
			})
			return
		}
		records = append(records, h)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Error reading database rows.",
			Error:   err.Error(),
		})
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    records,
	})
}

func (h *DbPool) CreateHealthData(c *gin.Context) {
	var req HealthDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid JSON body",
			Error:   err.Error(),
		})
	}

	ctx := c.Request.Context()

	query := `
        INSERT INTO MavHealth (date, mentalstate, physicalstate, note)
        VALUES ($1, $2, $3, $4)
        RETURNING id;`

	var newID int

	err := h.db.QueryRow(ctx, query, time.Now(), req.MentalState, req.PhysicalState).Scan(&newID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to insert health data",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    gin.H{"id": newID},
	})
}
