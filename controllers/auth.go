package controllers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Zavulon39/online-shop/schemas"
	"github.com/Zavulon39/online-shop/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (c *AuthController) Registration(ctx *gin.Context, db *sql.DB) {
	var regSchema schemas.RegistrationSchema
	var userId int

	if err := ctx.ShouldBindJSON(&regSchema); err != nil {
		ctx.JSON(422, gin.H{"detail": "Invalid data"})
		return
	}

	hashedPassword := services.GetMD5Hash(regSchema.Password)

	_, err := db.Exec("INSERT INTO users(login, password) VALUES($1, $2)", regSchema.Login, hashedPassword)

	if err != nil {
		ctx.JSON(500, gin.H{"detail": fmt.Sprintf("Error, while creating a user %v", err)})
		return
	}

	db.QueryRow("SELECT currval('users_id_seq')").Scan(&userId)

	_, err = db.Exec("INSERT INTO baskets(user_id) VALUES($1)", userId)

	if err != nil {
		ctx.JSON(500, gin.H{"detail": fmt.Sprintf("Error, while creating a user %v", err)})
		return
	}

	accessToken, accessErr := services.CreateAccessToken(userId)
	refreshToken, refreshErr := services.CreateRefreshToken(userId)

	if accessErr != nil || refreshErr != nil {
		ctx.JSON(500, gin.H{"detail": fmt.Sprintf("Error, while generating tokens: %v %v", accessErr, refreshErr)})
		return	
	}

	ctx.JSON(201, gin.H{"access_token": accessToken, "refresh_token": refreshToken, "username": regSchema.Login})
}

func (c *AuthController) Login(ctx *gin.Context, db *sql.DB) {
	var logingSchema schemas.RegistrationSchema
	var userId int

	if err := ctx.ShouldBindJSON(&logingSchema); err != nil {
		ctx.JSON(422, gin.H{"detail": "Invalid data"})
		return
	}

	hashedPassword := services.GetMD5Hash(logingSchema.Password)

	db.QueryRow("SELECT id FROM users WHERE login=$1 AND password=$2", logingSchema.Login, hashedPassword).Scan(&userId)

	if userId == 0 {
		ctx.JSON(404, gin.H{"detail": "User not found!"})
		return
	}

	accessToken, accessErr := services.CreateAccessToken(userId)
	refreshToken, refreshErr := services.CreateRefreshToken(userId)

	if accessErr != nil || refreshErr != nil {
		ctx.JSON(500, gin.H{"detail": fmt.Sprintf("Error, while generating tokens: %v %v", accessErr, refreshErr)})
		return	
	}

	ctx.JSON(200, gin.H{"access_token": accessToken, "refresh_token": refreshToken, "username": logingSchema.Login})
}

func (c *AuthController) RefreshToken(ctx *gin.Context, db *sql.DB) {
	var refreshSchema schemas.RefreshAccessTokenSchema
	var queryUserId int

	if err := ctx.ShouldBindJSON(&refreshSchema); err != nil {
		ctx.JSON(422, gin.H{"detail": "Invalid data!"})
		return
	}

	token, err := services.ParseToken(refreshSchema.RefreshToken)

	if err != nil {
		ctx.JSON(403, gin.H{"detail": fmt.Sprintf("%v", err)})
		return
	}
	if !token.Valid {
		ctx.JSON(403, gin.H{"detail": "Token signature is invalid!"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		ctx.JSON(403, gin.H{"detail": "Invalid refresh token!"})
		return
	}

	userId, err := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))

	db.QueryRow("SELECT id FROM users WHERE id=$1", userId).Scan(&queryUserId)

	if queryUserId == 0 || err != nil {
		ctx.JSON(404, gin.H{"detail": "Can't find user!"})
		return
	}

	accessToken, accessErr := services.CreateAccessToken(userId)

	if accessErr != nil {
		ctx.JSON(500, gin.H{"detail": fmt.Sprintf("Error, while generating tokens: %v", accessErr)})
		return	
	}

	ctx.JSON(200, gin.H{"access_token": accessToken, "refresh_token": refreshSchema.RefreshToken})
}