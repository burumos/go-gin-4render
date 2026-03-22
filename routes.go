package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}
