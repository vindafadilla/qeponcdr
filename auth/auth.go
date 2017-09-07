package ginAuth

// this package is a secure cookie authentication middleware for the the Gin Web Framework
// https://github.com/gin-gonic/gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"time"
	"strconv"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
)

const VERSION = "0.0.1"

// set our global package variables
var (
	CookieName string                       // the name of the cookie that will be used, default: "token"
	ConfigPath string                       // path to config file, default: ""
	ConfigType string                       // type of config file, default: "ini"
	Prefix string                           // the key in ctx.Keys[] to use, default: ""
	HashKey []byte                          // hash key for securecookie
	BlockKey []byte                         // block key for securecookie
	Expiration int64                        // time until the cookie expires in seconds, default: 604800
	Unauthorized func(ctx *gin.Context)     // function called if user is not authorized
	Authorized func(ctx *gin.Context)       // function called if user is authorized
	SecureCookie *securecookie.SecureCookie // global secure cookie object
)

func init() {

	ConfigPath = ""
	ConfigType = "ini"
	CookieName = "token"
	Prefix = ""
	Expiration = 604800 // 7 days

}

// gin middleware handler
// call this on your groups that require authentication
func Use(ctx *gin.Context) {

	err := Check(ctx)

	if err == nil {

		loggedIn, _ := ctx.Get(Prefix + "loggedIn")

		if loggedIn == true {

			if Authorized != nil {
				Authorized(ctx)
			}

		} else {

			if Unauthorized != nil {
				Unauthorized(ctx)
			}

		}

	}

}

// this function loads your specified configuration file and it's values
func LoadConfig() error {

		CookieName = "token"
		Prefix = ""

			val, err := hex.DecodeString("44adc9be4c44e848fb95a64e9fe5beef3c80215526e6b97d4399072dfb30951463b28c3336ebb31146d4ad13a8f04b0840ba5514303aa35a7aa30d2bc9ba7242")
			if err != nil {
				return err
			}

			HashKey = val
		
			val2, err := hex.DecodeString("268a231d4fe8d44ac79fc8cc7193215521e301e7e2bf80366c520de5246236fd")
			if err != nil {
				return err
			}

			BlockKey = val2

			val3, err := strconv.ParseInt("2600", 10, 64)
			if err != nil {
				return err
			}

			Expiration = val3
	

	return nil
}

// a private function that simply saves the log in status of the user to the current context
func saveLogin(ctx *gin.Context, status bool) {
	ctx.Set(Prefix+"loggedIn", status)
}

// a private function that returns the ip from the current request
func ip(ctx *gin.Context) string {
	return strings.Split(ctx.Request.RemoteAddr, ":")[0]
}

// checks for our token cookie, decodes it, and determines if it is valid
// the encrypted cookie data set in Login() will be set to the current context as well
func Check(ctx *gin.Context) error {

	// get the encrypted cookie value
	cookie, err := ctx.Request.Cookie(CookieName)

	if err == nil {

		data := make(map[string]string)

		SecureCookie = securecookie.New(HashKey, BlockKey)
		if err := SecureCookie.Decode(CookieName, cookie.Value, &data); err == nil {

			// save the login cookie data to the context
			ctx.Set(Prefix+"cookieData", data)

			hash := hashHeader(ctx)

			expiration, err := strconv.ParseInt(data["expiration"], 10, 64)
			if err != nil {
				return err
			}

			if hash == data["hash"] && ip(ctx) == data["ip"] && time.Now().Before(time.Unix(expiration, 0)) {

				saveLogin(ctx, true)

			} else {
				//call the full logout because it'll remove the cookie as well
				Logout(ctx)
			}

		} else {
			return err
		}

	}

	return nil
}

// handles the login process
// the first param is a map of strings that will be added to the cookie data before encryption and will be
// able to be recovered when Check() is called
func Login(ctx *gin.Context, extra map[string]string) error {

	data := make(map[string]string)

	for key, value := range extra {

		if key == "ip" || key == "hash" || key == "experation" {
			return errors.New("The key '" + key + "' is reserved.")
		}

		data[key] = value
	}

	// our current time + our expiration time, converted to a unix time stamp
	data["expiration"] = strconv.FormatInt(time.Now().Add(time.Duration(Expiration)*time.Second).Unix(), 10)
	data["ip"] = ip(ctx)
	data["hash"] = hashHeader(ctx)

	// encode our cookie data securely
	SecureCookie = securecookie.New(HashKey, BlockKey)
	if encoded, err := SecureCookie.Encode(CookieName, data); err == nil {

		//set our cookie
		cookie := http.Cookie{Name: CookieName, Value: encoded, Path: "/", MaxAge: int(Expiration)}
		http.SetCookie(ctx.Writer, &cookie)

	} else {
		return err
	}

	return nil
}

// removes our token cookie, sets the context to: not logged in
func Logout(ctx *gin.Context) {

	cookie := http.Cookie{Name: CookieName, Path: "/", MaxAge: -1}
	http.SetCookie(ctx.Writer, &cookie)
	saveLogin(ctx, false)

}

// this function returns and md5 hash (string) for a few common request headers
func hashHeader(ctx *gin.Context) string {

	h := md5.New()

	io.WriteString(h, ctx.Request.Header.Get("User-Agent"))
	io.WriteString(h, ctx.Request.Header.Get("Accept-Language"))

	return hex.EncodeToString(h.Sum(nil))

}
