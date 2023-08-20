package main

import (
	//เรียกใช้งาน func จากไฟล์ /controller/auth
	AuthController "TurterDev/go-jwt-api/controller/auth"
	UserController "TurterDev/go-jwt-api/controller/user"
	"TurterDev/go-jwt-api/middleware"
	"log"

	//เรียกใช้งาน db จากโฟเดอร์ /orm
	"TurterDev/go-jwt-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}
	orm.InitDB()

	r := gin.Default()
	// CORS allows all origins
	r.Use(cors.Default())
	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)
	authorized := r.Group("/users", middleware.JWTAuthen())
	authorized.GET("/readall", UserController.ReadAll)
	authorized.GET("/profile", UserController.Profile)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
