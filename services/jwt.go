package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateAccessToken(userId int) (string, error) {
  var err error

  atClaims := jwt.MapClaims{}

  atClaims["user_id"] = userId
  atClaims["type"] = "access"
  atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

  at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

  token, err := at.SignedString([]byte("melting_bananas_324"))
	
  if err != nil {
     return "", err
  }

  return token, nil
}

func CreateRefreshToken(userId int) (string, error) {
  var err error

  atClaims := jwt.MapClaims{}

  atClaims["user_id"] = userId
  atClaims["type"] = "refresh"
  atClaims["exp"] = time.Now().Add(time.Minute * 10080).Unix()

  at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

  token, err := at.SignedString([]byte("melting_bananas_324"))
	
  if err != nil {
     return "", err
  }

  return token, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
  return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("unexpexted singing method %v", token.Header)
    }

    return []byte("melting_bananas_324"), nil
  })
}