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
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *Application) AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("id")

		if userQueryId == "" {
			c.Header("Content-Type", "application/json")
			log.Println("user id is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty!"})
			c.Abort()
			return
		}
		user_id, err := primitive.ObjectIDFromHex(userQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err)
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: user_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "internal Server Error")
		}

		var addressInfo []bson.M
		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			log.Println(err)
			panic(err)
		}

		var size int32

		for _, json := range addressInfo {
			count := json["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: user_id}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

			if _, err = UserCollection.UpdateOne(ctx, filter, update); err != nil {
				log.Println(err)
			}
		} else {
			c.IndentedJSON(400, "Not Allowed to add more than 2 addresses.")
		}

		defer cancel()
		ctx.Done()

	}
}

func (app *Application) EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("id")

		if userQueryId == "" {
			c.Header("Content-Type", "application/json")
			log.Println("user id is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty!"})
			c.Abort()
			return
		}
		user_id, err := primitive.ObjectIDFromHex(userQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var editAddress models.Address

		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editAddress.House}, {Key: "address.0.street_name", Value: editAddress.Street}, {Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.PinCode}}}}
		if _, err = UserCollection.UpdateOne(ctx, filter, update); err != nil {
			c.IndentedJSON(404, "something went wrong trying to update")
			return
		}

		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the home address")
	}
}

func (app *Application) EditWorkHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("id")

		if userQueryId == "" {
			c.Header("Content-Type", "application/json")
			log.Println("user id is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty!"})
			c.Abort()
			return
		}
		user_id, err := primitive.ObjectIDFromHex(userQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: user_id}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editAddress.House}, {Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City}, {Key: "address.1.pin_code", Value: editAddress.PinCode}}}}
		if _, err = UserCollection.UpdateOne(ctx, filter, update); err != nil {
			c.IndentedJSON(404, "something went wrong trying to update")
			return
		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the work address")
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
		defer cancel()
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
