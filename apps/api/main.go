package main

import (
	"fmt"
	"github.com/123508/douyinshop/apps/api/handlers/address"
	"github.com/123508/douyinshop/apps/api/handlers/ai"
	"github.com/123508/douyinshop/apps/api/handlers/cart"
	"github.com/123508/douyinshop/apps/api/handlers/image"
	"github.com/123508/douyinshop/apps/api/handlers/order"
	"github.com/123508/douyinshop/apps/api/handlers/payment"
	"github.com/123508/douyinshop/apps/api/handlers/product"
	"github.com/123508/douyinshop/apps/api/handlers/shop"
	"github.com/123508/douyinshop/apps/api/handlers/user"
	"github.com/123508/douyinshop/apps/api/middleware"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// 启动 Prometheus metrics 服务器
	go func() {
		defer wg.Done()
		metricsAddr := fmt.Sprintf("%s:%d", config.Conf.HertzConfig.Host, 10000)
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Prometheus metrics server starting on %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			log.Printf("Prometheus metrics server error: %v", err)
		}
	}()

	// 启动主服务
	go func() {
		defer wg.Done()
		hertzAddr := fmt.Sprintf("%s:%d", config.Conf.HertzConfig.Host, config.Conf.HertzConfig.Port)
		hz := server.New(server.WithHostPorts(hertzAddr))

		// 注册 Prometheus 中间件
		middleware.RegisterPrometheus(hz)

		userGrop := hz.Group("/user")
		userGrop.POST("/register", user.Register)
		userGrop.POST("/login", user.Login)
		userGrop.GET("/logout", user.Logout)
		userGrop.GET("/info", middleware.ParseToken(), user.GetInfo)
		userGrop.POST("/update_info", middleware.ParseToken(), user.UpdateInfo)
		userGrop.DELETE("/delete", middleware.ParseToken(), user.Delete)

		productGroup := hz.Group("/product")
		productGroup.Use(middleware.ParseToken())
		productGroup.GET("/list", product.List)
		productGroup.GET("/detail/:product_id", product.Detail)
		productGroup.GET("/search", product.Search)

		cartGroup := hz.Group("/cart")
		cartGroup.Use(middleware.ParseToken())
		cartGroup.GET("/list", cart.List)
		cartGroup.POST("/add", cart.Add)
		cartGroup.DELETE("/empty", cart.Empty)
		cartGroup.POST("/reduce", cart.Reduce)
		cartGroup.POST("/checkout", cart.Checkout)

		shopGroup := hz.Group("/shop")
		shopGroup.Use(middleware.ParseToken())
		shopGroup.GET("/info/:shop_id", shop.GetInfo)
		shopGroup.GET("/getShopId", shop.GetShopId)
		shopGroup.GET("/list", shop.List)
		shopGroup.POST("/register", shop.Register)
		shopGroup.POST("/add", shop.Add)
		shopGroup.POST("/updateShop", shop.UpdateShopInfo)
		shopGroup.DELETE("/delete", shop.Delete)
		shopGroup.POST("/updateProduct", shop.UpdateProductInfo)

		addressGroup := hz.Group("/address")
		addressGroup.Use(middleware.ParseToken())
		addressGroup.GET("/list", address.List)
		addressGroup.POST("/delete", address.Delete)
		addressGroup.POST("/add", address.Add)
		addressGroup.POST("/update", address.Update)
		addressGroup.POST("/setDefaultAddr", address.SetDefault)

		paymentGroup := hz.Group("/payment")
		paymentGroup.Use(middleware.ParseToken())
		paymentGroup.POST("/charge", payment.Charge)
		paymentGroup.POST("/notify", payment.Notify)

		aiGroup := hz.Group("/ai")
		aiGroup.Use(middleware.ParseToken())
		aiGroup.GET("/orderQuery/:order_id", ai.OrderQuery)
		aiGroup.POST("/autoPlaceOrder", ai.AutoPlaceOrder)

		orderGroup := hz.Group("/order")
		orderGroup.Use(middleware.ParseToken())
		userOrderGroup := orderGroup.Group("/user")
		userOrderGroup.GET("/detail/:order_id", order.Detail)
		userOrderGroup.POST("/history", order.History)
		userOrderGroup.POST("/submit", order.Submit)
		userOrderGroup.POST("/cancel", order.Cancel)
		userOrderGroup.POST("/reminder", order.Reminder)
		userOrderGroup.GET("/complete/:order_id", order.Complete)

		shopOrderGroup := orderGroup.Group("/shop")
		shopOrderGroup.Use(middleware.ParseToken())
		shopOrderGroup.GET("/list", order.List)
		shopOrderGroup.GET("/detail/:shop_id", order.DetailShop)
		shopOrderGroup.GET("/confirm", order.Confirm)
		shopOrderGroup.GET("/delivery/:order_id", order.Delivery)
		shopOrderGroup.GET("/receive/:order_id", order.Receive)
		shopOrderGroup.POST("/rejection", order.Rejection)
		shopOrderGroup.POST("/cancel", order.CancelShop)

		imageGroup := hz.Group("/image")
		imageGroup.POST("/upload", image.UploadImage)
		imageGroup.GET("/get/:image", image.GetImage)

		log.Printf("Main server starting on %s", hertzAddr)
		if err := hz.Run(); err != nil {
			log.Printf("Main server error: %v", err)
		}
	}()

	// 等待所有服务启动
	wg.Wait()
}
