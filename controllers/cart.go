package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/MrBooi/ecommerce-cart/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	ProdCollection *mongo.Collection
	UserCollection *mongo.Collection
}

func NewApplication(prodCollection, UserCollection *mongo.Collection) *Application {
	return &Application{
		ProdCollection: prodCollection,
		UserCollection: UserCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		userQueryId := c.Query("userID")

		if productQueryId == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty!"))
			return
		}

		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty!"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.AddProductToCart(ctx, app.ProdCollection, app.UserCollection, productID, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully added to cart.")

	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {

	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		userQueryId := c.Query("userID")

		if productQueryId == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty!"))
			return
		}

		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty!"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.RemoveCartItem(ctx, app.ProdCollection, app.UserCollection, productID, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully remove item to cart.")

	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("userID")

		if userQueryId == "" {
			log.Panicln("user is is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty!"))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		err := database.BuyItemFromCart(ctx, app.UserCollection, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully bought item.")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		userQueryId := c.Query("userID")

		if productQueryId == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty!"))
			return
		}

		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty!"))
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.InstantBuy(ctx, app.ProdCollection, app.UserCollection, productID, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully placed the order.")
	}
}
