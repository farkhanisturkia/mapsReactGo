package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()

    // Enable CORS agar bisa diakses dari Vite dev server (http://localhost:5173)
    r.Use(cors.Default())

    r.GET("/api/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello from Go Backend!",
        })
    })

    r.Run(":9990")
}
