package middlewares

import (
	"fmt"

	"github.com/Zavulon39/online-shop/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var REQUIRED_AUTH_URLS = [3]string{
	"/api/add-to-basket/",
	"/api/remove-from-basket/",
	"/api/basket/",
}

func JWTLoginRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if Contains(ctx.Request.URL.Path) {
			const TOKEN_PREFIX = "Bearer "

			authHeader := ctx.GetHeader("Authorization")
			stringToken := authHeader[len(TOKEN_PREFIX):]
			token, err := services.ParseToken(stringToken)

			if err != nil {
				ctx.AbortWithStatusJSON(403, gin.H{"detail": fmt.Sprintf("%v", err)})
				return
			}
			if !token.Valid {
				ctx.AbortWithStatusJSON(403, gin.H{"detail": "Token signature is invalid!"})
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			// parseStr := strings.Split(strings.Replace(fmt.Sprintf("%v", claims["exp"]), ".", "", 1), "e")[0]
			// exp, _ := strconv.ParseInt(parseStr, 10, 64)
			// if exp < time.Now().Unix() {
			// 	ctx.AbortWithStatusJSON(403, gin.H{"detail": "Access token is expired!"})
			// 	return
			// }

			if claims["type"] != "access" {
				ctx.AbortWithStatusJSON(403, gin.H{"detail": "Invalid access token!"})
				return
			}
		}
	}
}

func Contains(url string) bool {
	for _, path := range REQUIRED_AUTH_URLS {
			if path == url {
					return true
			}
	}
	return false
}