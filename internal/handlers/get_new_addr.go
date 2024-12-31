package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type PostReqData struct {
	Currency string `json:"currency"`
}

func (h *Handler) GetNewAddr(c *gin.Context) {
	var req PostReqData

	log.Println("Received POST request")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal(err)
	}

	log.Printf("Decoded request: %+v", req)

	address, err := h.Runner.Ucase.GenerateNewAdd(req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		log.Fatal(err)
	}

	//SAIVING DATA TO GIN CONTEXT

	c.Set("generated_address", address)
	c.Set("currency", req.Currency)

	c.JSON(http.StatusOK, gin.H{
		"message": "Address created successfully",
		"address": address,
	})
}
