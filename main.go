package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres" 
    "gorm.io/gorm"
)

type Task struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
}

var db *gorm.DB

func main() {
    var err error

   
    dsn := "user=postgres dbname=textile password=1234 host=localhost port=5432 sslmode=disable TimeZone=UTC"

 
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }
    db.AutoMigrate(&Task{})

    r := gin.Default()

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"*"}
    config.AllowHeaders = []string{"Authorization", "Content-Type"}
    r.Use(cors.New(config))

    r.GET("/tasks", getTasks)
    r.POST("/tasks", createTask)
    r.GET("/tasks/:id", getTask)
    r.PUT("/tasks/:id", updateTask)
    r.DELETE("/tasks/:id", deleteTask)

    r.Run(":8080")
}


func getTasks(c *gin.Context) {
    var tasks []Task
    if err := db.Find(&tasks).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Görevleri getirme hatası"})
        return
    }

    c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
    var input struct {
        Title string `json:"title"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    task := Task{
        Title:     input.Title,
        Completed: false,
        CreatedAt: time.Now(),
    }

    if err := db.Create(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Görev oluşturma hatası"})
        return
    }

    c.JSON(http.StatusCreated, task)
}

func getTask(c *gin.Context) {
    id := c.Param("id")

    var task Task
    if err := db.Where("id = ?", id).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Görev bulunamadı"})
        return
    }

    c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
    id := c.Param("id")

    var input struct {
        Title     string `json:"title"`
        Completed bool   `json:"completed"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var task Task
    if err := db.Where("id = ?", id).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Görev bulunamadı"})
        return
    }

    task.Title = input.Title
    task.Completed = input.Completed

    if err := db.Save(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Görev güncelleme hatası"})
        return
    }

    c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
    id := c.Param("id")

    var task Task
    if err := db.Where("id = ?", id).First(&task).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Görev bulunamadı"})
        return
    }

    if err := db.Delete(&task).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Görev silme hatası"})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}
