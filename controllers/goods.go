package controllers

import (
	"database/sql"

	"github.com/Zavulon39/online-shop/schemas"
	"github.com/Zavulon39/online-shop/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type GoodsController struct{}

func (c *GoodsController) GetAllGoods(ctx *gin.Context, db *sql.DB) {
	goods := []schemas.Good{}

	rows, err := db.Query("SELECT * FROM goods")

	if err != nil {
		ctx.JSON(500, gin.H{"detail": "Can not get all goods!"})
		return
	}
	
	for rows.Next(){ 
    g := schemas.Good{}

    rows.Scan(&g.ID, &g.Title, &g.Description, &g.Price, &g.Count)

    goods = append(goods, g)
	}

	rows.Close()

	ctx.IndentedJSON(200, goods)
}

func (c *GoodsController) GetGoodsInBasket(ctx *gin.Context, db *sql.DB) {
	var basketID int
	var goods []schemas.Good

	const TOKEN_PREFIX = "Bearer "

	authHeader := ctx.GetHeader("Authorization")
	stringToken := authHeader[len(TOKEN_PREFIX):]
	token, _ := services.ParseToken(stringToken)

	userId := token.Claims.(jwt.MapClaims)["user_id"]

	db.QueryRow("SELECT id FROM baskets WHERE user_id=$1", userId).Scan(&basketID)

	if basketID == 0 {
		ctx.JSON(500, gin.H{"detail": "Can not find basket for user!"})
		return
	}

	rows, err := db.Query("SELECT * FROM goods WHERE id IN (SELECT good_id FROM baskets_goods WHERE basket_id=$1)", basketID)

	if err != nil {
		ctx.JSON(500, gin.H{"detail": "Error, while getting goods!"})
		return
	}

	for rows.Next(){
		good := schemas.Good{}
    rows.Scan(&good.ID, &good.Title, &good.Description, &good.Price, &good.Count)

		goods = append(goods, good)
	}

	rows.Close()

	ctx.IndentedJSON(200, goods)
}

func (c *GoodsController) AddGoodToBasket(ctx *gin.Context, db *sql.DB) {
	var addSchema schemas.AddGoodSchema
	var goodID int
	var basketID int

	if err := ctx.ShouldBindJSON(&addSchema); err != nil {
		ctx.JSON(400, gin.H{"detail": "Invalid data!"})
		return
	}

	db.QueryRow("SELECT id FROM goods WHERE id=$1", addSchema.GoodID).Scan(&goodID)

	if goodID == 0 {
		ctx.JSON(404, gin.H{"detail": "Can not find good!"})
		return
	}

	const TOKEN_PREFIX = "Bearer "

	authHeader := ctx.GetHeader("Authorization")
	stringToken := authHeader[len(TOKEN_PREFIX):]
	token, _ := services.ParseToken(stringToken)

	userId := token.Claims.(jwt.MapClaims)["user_id"]

	db.QueryRow("SELECT id FROM baskets WHERE user_id=$1", userId).Scan(&basketID)

	if basketID == 0 {
		ctx.JSON(500, gin.H{"detail": "Can not find basket for user!"})
		return
	}

	db.Exec("INSERT INTO baskets_goods VALUES($1, $2)", goodID, basketID)

	ctx.JSON(201, gin.H{})
}

func (c *GoodsController) RemoveGoodFromBasket(ctx *gin.Context, db *sql.DB) {
	var rmSchema schemas.AddGoodSchema
	var goodID int
	var basketID int

	if err := ctx.ShouldBindJSON(&rmSchema); err != nil {
		ctx.JSON(400, gin.H{"detail": "Invalid data!"})
		return
	}

	db.QueryRow("SELECT id FROM goods WHERE id=$1", rmSchema.GoodID).Scan(&goodID)

	if goodID == 0 {
		ctx.JSON(404, gin.H{"detail": "Can not find good!"})
		return
	}

	const TOKEN_PREFIX = "Bearer "

	authHeader := ctx.GetHeader("Authorization")
	stringToken := authHeader[len(TOKEN_PREFIX):]
	token, _ := services.ParseToken(stringToken)

	userId := token.Claims.(jwt.MapClaims)["user_id"]

	db.QueryRow("SELECT id FROM baskets WHERE user_id=$1", userId).Scan(&basketID)

	if basketID == 0 {
		ctx.JSON(500, gin.H{"detail": "Can not find basket for user!"})
		return
	}

	db.Exec("DELETE FROM baskets_goods WHERE good_id=$1 AND basket_id=$2", goodID, basketID)

	ctx.JSON(204, gin.H{})
}
