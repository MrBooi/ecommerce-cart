package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MrBooi/ecommerce-cart/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *Application) AddAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func (app *Application) EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (app *Application) EditWorkHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (app *Application) DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			log.Println("user id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)

		user_id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal Server Error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(404, "something went wrong trying to update")
			return
		}

		defer cancel()

		ctx.Done()
		c.IndentedJSON(200, "Successfully Deleted")

	}
}
