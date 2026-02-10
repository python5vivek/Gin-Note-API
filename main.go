package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"crypto/rand"
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomToken(n int) (string, error) {
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

type Loginuser struct {
	Username string
	Password string
}

func SignUp(c *gin.Context) {
	var user User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	db.Create(&user)
	tokenStr, _ := GenerateRandomToken(64)
	token := Token{UserId: user.ID, Token: tokenStr}
	db.Create(token)
	c.JSON(201, gin.H{"Token": token})
}

func Login(c *gin.Context) {
	var luser Loginuser
	err := c.ShouldBindJSON(luser)
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

}

func GetNote(c *gin.Context) {
	var notes []Note
	if err := db.Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func AddNote(c *gin.Context) {
	var not Note
	err := c.ShouldBindJSON(&not)
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	db.Create(&not)
	c.JSON(201, gin.H{"Message": "Created"})
}

var db *gorm.DB

func UpdateNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
	}
	var notes Note
	result := db.First(&notes, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"Message": "Note Not Found"})
		return
	}
	var updatedData Note
	err1 := c.ShouldBindJSON(&updatedData)
	if err1 != nil {
		c.JSON(400, gin.H{"Error": err1.Error()})
		return
	}
	db.Model(&notes).Updates(&updatedData)
	c.JSON(200, gin.H{"Message": "Updated"})
}

func DeleteNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	result := db.Delete(&Note{}, id)
	c.JSON(200, gin.H{"Total Deletes": result.RowsAffected})
}

func main() {
	db = InitDB()
	defer CloseDB(db)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/notes", GetNote)
	router.POST("/note", AddNote)
	router.PUT("/note/:id", UpdateNote)
	router.DELETE("/note/:id", DeleteNote)
	router.POST("/signup", SignUp)

	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
