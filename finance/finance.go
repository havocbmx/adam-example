package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type companyrevenue struct {
	CompanyId    string
	YTD          float64
	MTD          float64
	LastMonth    float64
	CurrencyCode string
}

func main() {
	bind := "0.0.0.0:9999"
	router := gin.New()
	apiKey := "abc123"

	router.Use(gin.Recovery())
	router.GET("/ping", func(c *gin.Context) { c.AbortWithStatusJSON(200, gin.H{"message": "pong"}) })
	router.Use(gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(string) bool { return true },
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Forwarded-For"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour, // - Preflight requests cached for 24 hours
	}))

	if apiKey != "" {
		router.Use(func(c *gin.Context) {
			if c.Query("Api-Key") == apiKey || c.GetHeader("Authorization") == apiKey {
				c.Next()
				return
			}
			c.AbortWithStatus(401)
		})
	}

	// router.GET("/tree/stats", func(c *gin.Context) {

	// 	c.JSON(200, pts.Stats)
	// 	// Start := time.Now()
	// 	// defer FireStatsReqHandler(c, 200, time.Since(Start).Milliseconds())
	// })

	router.GET("/finance/company/:id/revenue", func(c *gin.Context) {
		result := companyrevenue{}
		companyID := c.Param("id")
		if companyID != "" {
			result.CompanyId = companyID
			result.CurrencyCode = "USD"
			result.LastMonth = 1000.99
			result.MTD = 500
			result.YTD = 12008
			c.JSON(200, result)
		} else {
			c.JSON(404, "WRONG!!")

		}
	})

	log.Printf("Binding Finance API to: %s", bind)
	router.Run(bind)

}
