package database

import (
	"context"
	"errors"
	"log"
	"time"

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
	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Println(err)
		return ErrUserIIsNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"$usercart": bson.M{"_id": productId}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		return ErrCantRemoveItemCart
	}

	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userId string) error {
	// fetch the cart of the user
	// find the total of a cart
	//  create an order with the items
	// added order to the user collection
	//	empty up the cart
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_AT = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "usercart.price"}}}}}}
	pointerCursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getUsercart []bson.M

	if err = pointerCursor.All(ctx, &getUsercart); err != nil {
		panic(err)
	}

	var totalPrice int32

	for _, json := range getUsercart {
		price := json["total"]
		totalPrice = price.(int32)
	}

	orderCart.Price = int(totalPrice)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filterTwo := bson.D{primitive.E{Key: "_id", Value: id}}
	updateTwo := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}

	_, err = userCollection.UpdateOne(ctx, filterTwo, updateTwo)

	if err != nil {
		log.Println(err)
	}

	userCartEmpty := make([]models.ProductUser, 0)

	filterThree := bson.D{primitive.E{Key: "_id", Value: id}}
	updateThree := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: userCartEmpty}}}}
	_, err = userCollection.UpdateOne(ctx, filterThree, updateThree)

	if err != nil {
		return ErrCantBuyCartItem
	}

	return nil
}

func InstantBuy(ctx context.Context, productCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrUserIIsNotValid
	}

	var productDetails models.ProductUser
	var ordersDetail models.Order

	ordersDetail.Order_ID = primitive.NewObjectID()
	ordersDetail.Ordered_AT = time.Now()
	ordersDetail.Order_Cart = make([]models.ProductUser, 0)
	ordersDetail.Payment_Method.COD = true

	err = productCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productId}}).Decode(&productDetails)

	if err != nil {
		log.Println(err)
	}

	ordersDetail.Price = productDetails.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordersDetail}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println(err)
	}

	filterTwo := bson.D{primitive.E{Key: "_id", Value: id}}
	updateTwo := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}

	_, err = userCollection.UpdateOne(ctx, filterTwo, updateTwo)

	if err != nil {
		log.Println(err)
	}

	return nil
}
