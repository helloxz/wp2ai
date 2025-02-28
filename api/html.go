package api

import (
	"net/http"
	"wp2ai/utils"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"version": utils.Version,
		"date":    utils.VersionDate,
	})
}
