package handlers

import (
	"api_crypto1.0/internal/db/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetTxList(c *gin.Context) {
	var params repository.Params

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 4
	}
	data, err := h.Runner.Ucase.Repository.GetTxFromDB(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
