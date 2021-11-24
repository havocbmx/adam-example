package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Employee struct {
	EmployeeID string
	CompanyId  string
	FirstName  string
	LastName   string
	Active     bool
}

func main() {
	bind := "0.0.0.0:9998"
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

	router.GET("/hr/company/:id/employees", func(c *gin.Context) {

		companyID := c.Param("id")
		if companyID != "" {
			employees := getALLEmployees(companyID)
			c.JSON(200, employees)
		} else {
			c.JSON(404, "WRONG!!")

		}
	})

	router.GET("/hr/company/:id/employees/active", func(c *gin.Context) {

		companyID := c.Param("id")
		if companyID != "" {
			employees := getActiveEmployees(companyID)
			c.JSON(200, employees)
		} else {
			c.JSON(404, "WRONG!!")

		}
	})

	router.GET("/hr/company/:id/employees/inactive", func(c *gin.Context) {

		companyID := c.Param("id")
		if companyID != "" {
			employees := getInActiveEmployees(companyID)
			c.JSON(200, employees)
		} else {
			c.JSON(404, "WRONG!!")

		}
	})

	log.Printf("Binding Finance API to: %s", bind)
	router.Run(bind)

}

func getALLEmployees(companyID string) []Employee {

	employees := make([]Employee, 0)
	for i := 0; i < 30; i++ {
		employees = append(employees, getEmployee(i, companyID))
	}
	return employees
}

func getActiveEmployees(companyID string) []Employee {

	employees := make([]Employee, 0)
	for i := 0; i < 30; i++ {
		employee := getEmployee(i, companyID)
		if employee.Active {
			employees = append(employees, employee)
		}
	}
	return employees
}

func getInActiveEmployees(companyID string) []Employee {

	employees := make([]Employee, 0)
	for i := 0; i < 30; i++ {
		employee := getEmployee(i, companyID)
		if !employee.Active {
			employees = append(employees, employee)
		}
	}
	return employees
}

func getEmployee(employeeNumber int, companyid string) Employee {
	NumStr := strconv.Itoa(employeeNumber)
	var Active bool
	if employeeNumber%2 == 1 {
		Active = true
	}

	return Employee{EmployeeID: NumStr, FirstName: "EMP" + NumStr, LastName: "EMPLAST" + NumStr, Active: Active, CompanyId: companyid}
}
