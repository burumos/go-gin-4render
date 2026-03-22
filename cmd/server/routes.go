package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func dbVersion(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not connected"})
		return
	}
	var version string
	err := db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"version": version})
}

func formGET(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func formPOST(c *gin.Context) {
	data := c.PostForm("data")
	if data == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"Message": "Data is required"})
		return
	}

	// Insert or update (simple: always insert for demo)
	_, err := db.Exec("INSERT INTO test (data) VALUES ($1)", data)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"Message": "Error saving data"})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"Message": "Data saved successfully", "Data": data})
}

func loginGET(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginPOST(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"Error": "Username and password required"})
		return
	}

	var hash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&hash)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Invalid credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set("user", username)
	session.Save()

	c.Redirect(http.StatusFound, "/form")
}

func logoutGET(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}
