package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	log.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
