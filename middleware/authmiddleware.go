package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("secret"))

func Authenticate(c *gin.Context) bool {
	session, _ := Store.Get(c.Request, "jwt_token")
	token, ok := session.Values["token"]
	fmt.Println(token)
	if !ok {

		return ok
	}
	return true
}
