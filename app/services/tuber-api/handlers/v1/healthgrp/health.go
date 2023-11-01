package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func New() *HealthController {
	return &HealthController{}
}

func (h HealthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
