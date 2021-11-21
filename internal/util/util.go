package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AbortResponseUnexpectedError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred."})
}

func AbortResponseNotFound(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "This resource was not found."})
}
