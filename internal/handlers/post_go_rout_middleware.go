package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//USING THIS MIDLWARE TO START LISTENING TX AS SOON AS WE GOT NEW CREATED ADDRESS

func (h *Handler) PostRoutineMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Request.Method == http.MethodPost && c.FullPath() == "/create/address" {

			// GETTING DATA FROM CONTEXT
			address, exists := c.Get("generated_address")
			if !exists {
				log.Println("Generated address not found in context")
				return
			}
			currency, exists := c.Get("currency")
			if !exists {
				log.Println("Currency not found in context")
				return
			}

			addr := address.(string)
			curr := currency.(string)

			// STARTING LISTENING TX

			go h.Runner.BLockListener(addr, curr)
		}
	}
}
