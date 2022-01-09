package main

import (
	"io/ioutil"

	"github.com/Zavulon39/online-shop/controllers"
	"github.com/Zavulon39/online-shop/middlewares"
	"github.com/Zavulon39/online-shop/services"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	dbService := services.NewDBService()
	sql, _ := ioutil.ReadFile("init.sql")
	db := dbService.GetDB()

	server.Use(middlewares.JWTLoginRequired())

	dbService.InitSQL(string(sql))

	authController := controllers.AuthController{}
	goodsController := controllers.GoodsController{}

	api := server.Group("/api")
	{
		api.GET("/goods/", func(c *gin.Context) {
			goodsController.GetAllGoods(c, db)
		})
		api.GET("/basket/", func(c *gin.Context) {
			goodsController.GetGoodsInBasket(c, db)
		})

		api.POST("/add-to-basket/", func(c *gin.Context) {
			goodsController.AddGoodToBasket(c, db)
		})
		api.POST("/remove-from-basket/", func(c *gin.Context) {
			goodsController.RemoveGoodFromBasket(c, db)
		})

		auth := api.Group("/auth")
		{
			auth.POST("/login/", func(c *gin.Context) {
				authController.Login(c, db)
			})
			auth.POST("/registration/", func(c *gin.Context) {
				authController.Registration(c, db)
			})
			auth.POST("/refresh-access-token/", func(c *gin.Context) {
				authController.RefreshToken(c, db)
			})
		}
	}

	server.Run(":8000")
}