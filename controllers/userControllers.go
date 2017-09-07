package controllers

import (
	"github.com/vindafadilla/qeponcdr/qeponcdr/models"
	auth "github.com/vindafadilla/qeponcdr/qeponcdr/auth"
	
	"github.com/gorilla/securecookie"
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"

	"fmt"
	"time"
)

//---------------------------Methods and Variable for Database-----------------
var SESSION *mgo.Session

func init() {
    session, err := mgo.Dial("127.0.0.1")
    Must(err)
    SESSION = session
}

func MapMongo(c *gin.Context) {
    s := SESSION.Clone()

    defer s.Close()

    c.Set("mongo", s.DB("greenphonetemp"))
    c.Next()
}

func Must(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func initializeDB(ctx *gin.Context, collection string)(*mgo.Collection) {
	//Initialize database MongoDB 
	db := ctx.MustGet("mongo").(*mgo.Database)
    usercoll:= db.C(collection)
    return usercoll
}

//------------------------Methods for Authentication User----------------------
func LoginPost(ctx *gin.Context) {
	//Initialize collection MongoDB
	usercoll:=initializeDB(ctx,"adminuser")

	// Initialize data structs
	user := models.Administrator{}
	data := map[string]string{}

	// Expecting login data in JSON form
	ctx.Bind(&data)
	var passmatch bool

	//Search data's user with email and password in database
	err := usercoll.Find(bson.M{"email": data["email"]}).One(&user)
	if err != nil {
		ctx.String(403, "User was not found!")
		return
	}else{
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
		if err!=nil {
			ctx.String(403, "Password not match!")
			passmatch = false
		}else{
		// ctx.JSON(200, "")
			passmatch = true
		}
	}

	if passmatch == true{
		// this data will be added to the cookie and available on decode
		extra := map[string]string{"email": user.Email}

		// log in the user
		err1 := auth.Login(ctx, extra)
		if err1 == nil {
			ctx.JSON(200, user)
		}
	}
	

}

func LoginAuthenticated(ctx *gin.Context) {

	// fetch our decrypted data that was set on our cookie
	// NOTE: if you set a prefix you need to add it before cookieData
	cookie, err := ctx.Get(auth.Prefix + "loggedIn")
	if err == false {
		// return the email we set with auth.Login
		fmt.Println("The current logged in user's email is: ", cookie)
	}

}

func LoginUnauthenticated(ctx *gin.Context) {

	ctx.String(401, "You are not logged in!")

	// ctx.Abort() means no other handlers will be ran after this one, meaning the route controller won't execute
	ctx.Abort()

}

//-------------------------Methods for Registration User-----------------------
func RegisterUser(ctx *gin.Context) {

	//Initialize collection MongoDB
	usercoll:=initializeDB(ctx,"adminuser")
    
	// Initialize data structs
	user := models.Administrator{}

	// Expecting register data in JSON form
	extra := map[string]string{}
	ctx.Bind(&extra)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(extra["password"]), bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    
    Password :=string(hashedPassword)
    CreatedAt := time.Now().UTC()

	//Search user if it was already exist
	err = usercoll.Find(bson.M{"email": extra["email"]}).One(&user)
	if err != nil {
		err = usercoll.Insert(&models.Administrator{extra["email"],Password, extra["firstname"], extra["lastname"], extra["address"],extra["country"],extra["phone"], CreatedAt})
		if err !=nil{
			ctx.String(400, "Cannot insert!",err)
		}else{
			ctx.JSON(201, models.Administrator{extra["email"],Password, extra["firstname"], extra["lastname"], extra["address"],extra["country"],extra["phone"], CreatedAt})
		}
	}else{
		ctx.String(409, "Account already exist!")
	}

}

//-------------------------Methods for Update Data User------------------------
func UpdateUser(ctx *gin.Context) {

	//Initialize collection MongoDB
	usercoll:=initializeDB(ctx,"adminuser")

    //Initialize variable and struct
    dataupdate := map[string]string{}
    ctx.Bind(&dataupdate)
    
    //Query to database
	cookie, err := ctx.Request.Cookie(auth.CookieName)
	if err==nil{
		data := make(map[string]string)

		SecureCookie := securecookie.New(auth.HashKey, auth.BlockKey)
		if err := SecureCookie.Decode(auth.CookieName, cookie.Value, &data); err == nil {
			data2:= data["email"]
		// ctx.JSON(200, data2)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dataupdate["password"]), bcrypt.DefaultCost)
		    if err != nil {
		        panic(err)
		    }
			dataupdate["password"] = string(hashedPassword)
			colQuerier := bson.M{"email": data2}
			var change bson.M
			if dataupdate["firstname"]!= "" ||dataupdate["lastname"]!= "" || dataupdate["address"]!= "" ||dataupdate["country"]!= "" ||dataupdate["phone"]!= "" || dataupdate["photo"]!= ""{
				change = bson.M{"$set": bson.M{"firstname": dataupdate["firstname"],"lastname": dataupdate["lastname"],"address": dataupdate["address"],"country":dataupdate["country"],"phone": dataupdate["phone"]}}
			}else if dataupdate["password"]!=""{
				change = bson.M{"$set": bson.M{"password":dataupdate["password"]}}
			}
			
			
			err = usercoll.Update(colQuerier, change)
			if err != nil {
				ctx.String(400, "Cannot insert!",err)
			}else{
				ctx.JSON(200, dataupdate)
			}
			
		}

	}else{
		 ctx.String(400, "Cannot update!",err)
	}

}
