package routes


import ( 
	"github.com/MrBooi/ecommerce-cart/controllers" 
	"github.com/gin-gonic/gin"
)

func userRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup")
	incomingRoutes.POST("/user/login")
	incomingRoutes.POST("/admin/addproduct")
	incomingRoutes.POST("/users/productview")
	incomingRoutes.POST("/user/search", )
}