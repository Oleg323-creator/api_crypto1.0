package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) Withdraw(c *gin.Context) {
	var req PostReqData

	log.Println("Received POST request")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal(err)
	}

	log.Printf("Decoded request: %+v", req)

	err := h.Runner.Ucase.Withdraw(req.Address, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		log.Fatal(err)
	}

	id, err := h.Runner.Ucase.GenerateWithdrawID()
	if err != nil {
		log.Fatalf("Error generating withdraw ID: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "You have successfully withdrawn money, your id:",
		"id":      id,
	})
}
