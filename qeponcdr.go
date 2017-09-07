package main

import (
	"github.com/vindafadilla/qeponcdr/qeponcdr/controllers"
	auth "github.com/vindafadilla/qeponcdr/qeponcdr/auth"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	
	"fmt"	
)


func main() {

	r := gin.Default()
	r.Use(func(c *gin.Context) {
               // Run this on all requests   
               // Should be moved to a proper middleware 
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Next()
    })

    r.OPTIONS("/*cors", func(c *gin.Context) {
              // Empty 200 response
    })

	apiv1 := r.Group("/v1")
	{
		//Method to initialize MongoDB
		apiv1.Use(controllers.MapMongo)
		//Method for authentication
		//This method use data binding to get value in parameter REST (ex : email & password), then, it will be stored on a varible.
		//So, that value will be check in database, if data wasn't exist, so value will return response 400, where user not found and give status loginunauthenticated
		//If data was exists, then, it will be return status ok, that will send data's user and give status loginauthenticated and save cookie.
		apiv1.POST("/login", controllers.LoginPost)
		//Method for registration
		apiv1.POST("/register", controllers.RegisterUser)
		//Method for logout and delete cookie
		apiv1.GET("/logout", func(ctx *gin.Context) {

				// logs out the user
				auth.Logout(ctx)

			})

		// set up and authentication group that uses middleware
		//Method for check if user authenticated or not, it will give status login authenticated only if user was authenticated
		authenticate := apiv1.Group("/")
		// tell it to use the ginAuth middleware on these routes
		authenticate.Use(auth.Use)
		authenticate.GET("/checklogin", func(ctx *gin.Context){

				// if not logged in, you will never reach this
				ctx.String(200, "You are logged in!")

			})
		//Method for get data about total user, total user verified, new user and new verified user.
		
		authenticate.PUT("/updateuser",func(ctx *gin.Context) {

				// get data
				controllers.UpdateUser(ctx)

			})
		authenticate.GET("/getdatapermonth",func(ctx *gin.Context) {

				// get data
				controllers.GetDataperMonth(ctx)

			})

		authenticate.GET("/getalldatapermonth",func(ctx *gin.Context) {

				// get data
				controllers.GetAllDataperMonth(ctx)

			})

		authenticate.GET("/getdataperday",func(ctx *gin.Context) {

				// get data
				controllers.GetDataperDay(ctx)

			})

	}

	// Unauthorized handler
	auth.Unauthorized = controllers.LoginUnauthenticated
	// Authorized handler
	auth.Authorized = controllers.LoginAuthenticated

	// Load config file, can skip this and set the values yourself
	err := auth.LoadConfig()

	if err != nil {

		// If there was an error loading our config file
		fmt.Println(err)

	} else {

		// Run our server if no errors
		r.Run("127.0.0.1:8080")

	}

}
