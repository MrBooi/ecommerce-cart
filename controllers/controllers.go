package controllers

import ()

import "github.com/gin-gonic/gin"

func HashPassword(password string) string {}

func VerifyPassword(userPassword, givenPassword string) (bool, string) {}

func SignUp() gin.HandlerFunc {

}

func Login() gin.HandlerFunc {

}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProductByQuery() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

}
