package auth

import (
	"TurterDev/go-jwt-api/orm"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/golang-jwt/jwt"
	_ "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

// ประกาศ Body ที่จะใช้รับข้อมูล
type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar"`
}

func Register(c *gin.Context) {
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check User
	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(200, gin.H{"status": "Error", "message": "Create User Exist"})
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	if err != nil {
		log.Fatal(err)
	}

	//สร้าง OBJ สำหรับ save data to DB
	user := orm.User{Username: json.Username, Password: string(encryptedPassword), Fullname: json.Fullname, Avatar: json.Avatar}
	//Create User
	orm.Db.Create(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Create User Success", "userId": user.ID})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "Error", "message": "Create User Failed"})
	}
}

// ประกาศ Body ที่จะใช้รับข้อมูล
type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Check User
	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "Error", "message": "User Dose Not Exist"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_Secret_Key"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExist.ID,
			"exp":    time.Now().Add(time.Minute * 1).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)
		fmt.Println(tokenString, err)
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Login Success", "token": tokenString})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "Error", "message": "Login Failed"})
	}
}
