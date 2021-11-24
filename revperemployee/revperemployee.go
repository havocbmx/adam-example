package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

type Employee struct {
	EmployeeID string
	CompanyId  string
	FirstName  string
	LastName   string
	Active     bool
}

var EmployeeCount int
var LastKnownRev float64

func main() {
	startAPI("0.0.0.0:9997")
}

func getActiveEmployeesForCompany(company, accessKey string) ([]Employee, error) {
	employees, err := getEmployeesForCompany(company, accessKey)
	activeEmployees := make([]Employee, 0)
	for _, emp := range employees {
		if emp.Active {
			activeEmployees = append(activeEmployees, emp)
		}

	}
	return activeEmployees, err
}

func getEmployeesForCompany(company, accessKey string) ([]Employee, error) {
	url := "http://127.0.0.1:9998/hr/company/" + company + "/employees?Api-Key=" + accessKey

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	Employees := make([]Employee, 0)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Employees, err
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		return Employees, getErr
	}
	if res.StatusCode != 200 {

		if res.StatusCode == 401 {
			return Employees, errors.New("GOT A FORBIDDEN USE ACCESS KEY")
		} else {
			return Employees, errors.New("GOT A StatusCode of " + strconv.Itoa(res.StatusCode))
		}
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return Employees, readErr
	}

	jsonErr := json.Unmarshal(body, &Employees)
	if jsonErr != nil {
		return Employees, jsonErr
	}

	fmt.Println(Employees, "Count", len(Employees))

	return Employees, nil

}

func getRevenueForCompany(company, accessKey string) (float64, error) {

	url := "http://127.0.0.1:9999/finance/company/" + company + "/revenue?Api-Key=" + accessKey

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.StatusCode != 200 {

		if res.StatusCode == 401 {
			return 0, errors.New("GOT A FORBIDDEN USE ACCESS KEY")
		} else {
			return 0, errors.New("GOT A StatusCode of " + strconv.Itoa(res.StatusCode))
		}
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	companyrevresponse := companyrevenue{}
	jsonErr := json.Unmarshal(body, &companyrevresponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(companyrevresponse.LastMonth)

	return companyrevresponse.LastMonth, nil

}

func getRevEmployeeHandler(c *gin.Context) {
	companyID := c.Param("id")
	accessKey := "abc123"
	lastMonthRev, err := getRevenueForCompany(companyID, accessKey)

	if err != nil {
		fmt.Println(err)
		if err.Error() == "GOT A FORBIDDEN USE ACCESS KEY" {
			fmt.Println("Looks like we forgot our accesskey trying with key")
			accessKey := "abc123"
			lastMonthRev, err = getRevenueForCompany(companyID, accessKey)
			if err != nil {
				fmt.Println("Failed adding key for recovery exiting...")
				c.JSON(500, "Failed adding key for recovery exiting...")
				return
			}

		} else {
			fmt.Println("FAILED RECOVERY...EXITING")
			c.JSON(500, "FAILED RECOVERY...EXITING")
			return
		}
		//do error handling stuff
	}
	fmt.Println("I got ", lastMonthRev)
	var employeeCount int
	employees, empError := getActiveEmployeesForCompany(companyID, accessKey)
	if empError != nil {
		if EmployeeCount != 0 {
			fmt.Println(empError, "Recovered... Using Employee Count Cache", EmployeeCount)
			employeeCount = EmployeeCount
		} else {
			fmt.Println(empError, "exiting... Couldn't Use Cache")
			c.JSON(500, empError)
			return
		}

	} else {
		employeeCount = len(employees)
		EmployeeCount = employeeCount
	}

	if employeeCount > 0 {
		fmt.Println("rev per employee", lastMonthRev/float64(employeeCount))
		c.JSON(200, lastMonthRev/float64(employeeCount))
		return
	}
	c.JSON(200, 0)
}

func startAPI(bind string) {
	//bind := "0.0.0.0:9998"
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

	router.GET("/revemp/company/:id/emprev", getRevEmployeeHandler)

	// router.GET("/hr/company/:id/employees/active", func(c *gin.Context) {

	// 	companyID := c.Param("id")
	// 	if companyID != "" {
	// 		employees := getActiveEmployees(companyID)
	// 		c.JSON(200, employees)
	// 	} else {
	// 		c.JSON(404, "WRONG!!")

	// 	}
	// })

	// router.GET("/hr/company/:id/employees/inactive", func(c *gin.Context) {

	// 	companyID := c.Param("id")
	// 	if companyID != "" {
	// 		employees := getInActiveEmployees(companyID)
	// 		c.JSON(200, employees)
	// 	} else {
	// 		c.JSON(404, "WRONG!!")

	// 	}
	// })

	log.Printf("Binding Finance API to: %s", bind)
	router.Run(bind)
}
