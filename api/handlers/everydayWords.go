package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WordOftheday struct {
	ID   int    `json:"id"`
	Word string `json:"word"`
}

func (h *DbPool) GetWordOfTheDay(c *gin.Context) {
	ctx := c.Request.Context()
	var req WordOftheday
	err := h.db.QueryRow(ctx, "Select id, word from WordOfTheDay order by id desc limit 1").Scan(&req.ID, &req.Word)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "no words found womp",
				"success": false,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "internal server implosion",
		})
		return
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    req,
	})
}

func (h *DbPool) CreateWordOfTheDay(c *gin.Context) {
	var req WordOftheday
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"success": false,
			"message": "invalid request body",
		})
	}
	ctx := c.Request.Context()

	err := h.db.QueryRow(ctx, "insert into WordOfTheDay (word) values ($1);", req.Word)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to insert health data",
			Error:   "fuck",
		})
		return
	}
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
	})
}
