package router

import (
	"github.com/gin-gonic/gin"
	"softengineering/object"
	"softengineering/service"
)

var queue chan struct{}

func limit(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		select {
		case queue <- struct{}{}:
			defer func() {
				// Process end, remove a slot
				<-queue
			}()
			next(c)
		default:
			// If queue is full, wait here
			queue <- struct{}{}
			defer func() {
				<-queue
			}()
			next(c)
		}
	}
}
func InitRouter() {
	queue = make(chan struct{}, 3)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	router.LoadHTMLGlob("template/*")
	router.GET("/status", func(c *gin.Context) {
		c.HTML(200, "status.html",
			gin.H{
				"power": object.MasterACInstance.Power,
				"mode":  object.MasterACInstance.Mode,
				"temp":  object.MasterACInstance.Temperature,
			})
	})
	router.GET("/report", func(c *gin.Context) {
		c.HTML(200, "report.html", nil)
	})
	internal := router.Group("/internal-api")
	{
		internal.GET("/register", service.RegisterSlave)
		internal.POST("/modify", limit(service.ModifyAC))
		internal.POST("/serve", service.ServeSlave)
		internal.GET("/info", service.Info)
		internal.GET("/status", service.SlaveStatus)
		internal.GET("/fee", service.SlaveStatus)
		internal.GET("/logout", service.SlaveLogOut)
		internal.GET("/get-report", service.Report)
		internal.GET("/serve-status", service.PendingSlaveServe)
		internal.POST("/stop", service.StopSlaveServe)
		internal.GET("/edit", service.EditMaster)
	}
	router.Run(":8080")
}
