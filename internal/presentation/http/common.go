package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/maklybae/ddd-zoo/internal/types/openapi/v1"
)

func (s *Server) SendBadRequestResponse(c *gin.Context, err error, details map[string]interface{}) {
	response := v1.ApiErrorResponse{
		Details:   &details,
		Error:     err.Error(),
		Message:   "Something went wrong",
		Timestamp: s.timeProvider.Now(),
	}
	c.JSON(http.StatusBadRequest, response)
}
