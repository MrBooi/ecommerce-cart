package database

import (
	"context"
	"errors"
	"log"

	"github.com/MrBooi/ecommerce-cart/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product.")
	ErrCantDecodeProducts = errors.New("can't decode the products.")
	ErrUserIIsNotValid    = errors.New("this user is not valid.")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart.")
	ErrCantRemoveItemCart = errors.New("cannot remove this item from the cart.")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart.")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase.")
)

func AddProductToCart(ctx context.Context, productCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	searchFromDB, err := productCollection.Find(ctx, bson.M{"_id": productId})

	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser

	err = searchFromDB.All(ctx, &productCart)

	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Println(err)
		return ErrUserIIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, productCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {

	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userId string) error {
	return nil
}

func InstantBuy(ctx context.Context, productCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	return nil
}
