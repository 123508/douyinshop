package main

import (
	"github.com/123508/douyinshop/apps/api/handlers/address"
	"github.com/123508/douyinshop/apps/api/handlers/ai"
	"github.com/123508/douyinshop/apps/api/handlers/cart"
	"github.com/123508/douyinshop/apps/api/handlers/order"
	"github.com/123508/douyinshop/apps/api/handlers/payment"
	"github.com/123508/douyinshop/apps/api/handlers/product"
	"github.com/123508/douyinshop/apps/api/handlers/shop"
	"github.com/123508/douyinshop/apps/api/handlers/user"
	"github.com/123508/douyinshop/pkg/config"

	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	hertzAddr := fmt.Sprintf("%s:%d", config.Conf.HertzConfig.Host, config.Conf.HertzConfig.Port)
	hz := server.New(server.WithHostPorts(hertzAddr))

	userGrop := hz.Group("/user")
	userGrop.POST("/register", user.Register)
	userGrop.POST("/login", user.Login)
	userGrop.GET("/logout", user.Logout)
	userGrop.GET("/info", user.GetInfo)
	userGrop.POST("/update_info", user.UpdateInfo)
	userGrop.DELETE("/delete", user.Delete)

	productGroup := hz.Group("/product")
	productGroup.GET("/list", product.List)
	productGroup.GET("/detail/:product_id", product.Detail)
	productGroup.GET("/search", product.Search)

	cartGroup := hz.Group("/cart")
	cartGroup.GET("/list", cart.List)
	cartGroup.POST("/add", cart.Add)
	cartGroup.DELETE("/empty", cart.Empty)
	cartGroup.POST("/reduce", cart.Reduce)
	cartGroup.POST("/checkout", cart.Checkout)

	shopGroup := hz.Group("/shop")
	shopGroup.GET("/info/:shop_id", shop.GetInfo)
	shopGroup.GET("/getShopId", shop.GetShopId)
	shopGroup.GET("/list", shop.List)
	shopGroup.POST("/register", shop.Register)
	shopGroup.POST("/add", shop.Add)
	shopGroup.POST("/updateShop", shop.UpdateShopInfo)
	shopGroup.DELETE("/delete", shop.Delete)
	shopGroup.POST("/updateProduct", shop.UpdateProductInfo)

	addressGroup := hz.Group("/address")
	addressGroup.GET("/list", address.List)
	addressGroup.POST("/delete", address.Delete)
	addressGroup.POST("/add", address.Add)
	addressGroup.POST("/update", address.Update)
	addressGroup.POST("/setDefaultAddr", address.SetDefault)

	paymentGroup := hz.Group("/payment")
	paymentGroup.POST("/charge", payment.Charge)
	paymentGroup.POST("/notify", payment.Notify)

	aiGroup := hz.Group("/ai")
	aiGroup.GET("/orderQuery/:order_id", ai.OrderQuery)
	aiGroup.POST("/autoPlaceOrder", ai.AutoPlaceOrder)

	orderGroup := hz.Group("/order")
	userOrderGroup := orderGroup.Group("/user")
	userOrderGroup.GET("/detail/:order_id", order.Detail)
	userOrderGroup.POST("/history", order.History)
	userOrderGroup.POST("/submit", order.Submit)
	userOrderGroup.POST("/cancel", order.Cancel)
	userOrderGroup.POST("/reminder", order.Reminder)
	userOrderGroup.GET("/complete", order.Complete)

	shopOrderGroup := orderGroup.Group("/shop")
	shopOrderGroup.GET("/list", order.List)
	shopOrderGroup.GET("/detail/:shop_id", order.DetailShop)
	shopOrderGroup.GET("/confirm", order.Confirm)
	shopOrderGroup.GET("/delivery/:order_id", order.Delivery)
	shopOrderGroup.GET("/receive/:order_id", order.Receive)
	shopOrderGroup.POST("/rejection", order.Rejection)
	shopOrderGroup.POST("/cancel", order.CancelShop)

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}
