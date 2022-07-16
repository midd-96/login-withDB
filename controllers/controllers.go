package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"project_login/database"
	"project_login/helpers"
	"project_login/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store = sessions.NewCookieStore([]byte("secret"))

type Users struct {
	ID       int
	Username string
	Password string
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	if err != nil {
		check = false
	}
	return check

}

func Login(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := UserLoged(c)
	if ok {
		c.Redirect(303, "/home")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func PostLogin(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	// var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user []Users
	var status bool
	Fusername := c.Request.FormValue("username")
	Fpassword := c.Request.FormValue("password")

	db := database.InitDB()
	db.Find(&user)
	for _, i := range user {
		passwordIsValid := VerifyPassword(Fpassword, i.Password)
		if i.Username == Fusername && passwordIsValid {
			status = true
			break
		}
	}

	if !status {
		log.Println("Wrong Username or Password\n\t\tTry Again")
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	token, _, _ := helpers.GenerateTokens(Fusername, "User")
	session, _ := Store.Get(c.Request, "jwt_token")
	session.Values["token"] = token
	session.Save(c.Request, c.Writer)
	// defer Close()
	c.Redirect(http.StatusSeeOther, "/home")
}

func Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func PostSignup(c *gin.Context) {
	var user []Users
	var status bool = true
	FusernameN := c.Request.FormValue("username")
	Fpassword := HashPassword(c.Request.FormValue("password"))

	db := database.InitDB()
	db.AutoMigrate(&Users{})
	db.Find(&user)

	for _, i := range user {
		if i.Username == FusernameN {
			status = false
			break
		}
	}

	if !status {
		log.Printf("hello %s , The username is already taken", FusernameN)
		c.Redirect(303, "/signup")
		return

	}

	db.Create(&Users{Username: FusernameN, Password: Fpassword})
	log.Printf("Hey %s, Your account is successfully created.", FusernameN)
	c.Redirect(http.StatusSeeOther, "/login")

}

func Admin(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ok := AdminLoged(c)
	if ok {
		c.Redirect(303, "/wadmin")
		return
	}
	c.HTML(http.StatusOK, "admin.html", nil)
}

func PostAdmin(c *gin.Context) {

	config := &models.Admin{
		UserName: os.Getenv("ADMIN_NAME"),
		Password: HashPassword(os.Getenv("PASSWORD")),
	}

	config.Password = HashPassword(config.Password)

	Fusername := c.Request.FormValue("username")
	Fpassword := c.Request.FormValue("password")

	passwordIsValid := VerifyPassword(Fpassword, config.Password)

	if Fusername != config.UserName || passwordIsValid {
		log.Println("Wrong Username or Password , Check Again!")
		c.Redirect(303, "/admin")
		return
	}

	token, _, _ := helpers.GenerateTokens(Fusername, "Admin")

	session, _ := Store.Get(c.Request, "admin_jwt_token")
	session.Values["token"] = token
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusSeeOther, "/home")

	c.Redirect(303, "/wadmin")

}

func Wadmin(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	var user []Users

	ok := AdminLoged(c)
	if !ok {
		c.Redirect(303, "/admin")
		return
	}

	db := database.InitDB()
	var us = [11]string{}

	var id = [11]int{}
	db.Raw("SELECT id,username FROM users").Scan(&user)
	for ind, i := range user {
		us[ind], id[ind] = i.Username, i.ID

	}

	c.HTML(http.StatusOK, "welcomeadmin.html", gin.H{

		"users": us,
		"id":    id,
	})
}

func Home(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ok := UserLoged(c)
	if !ok {
		c.Redirect(303, "/login")
		return
	}
	c.HTML(http.StatusOK, "welcomeuser.html", nil)
}

func Logout(c *gin.Context) {

	cookie, err := c.Request.Cookie("jwt_token")
	if err != nil {
		c.Redirect(303, "/login")
	}
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, false)
	_ = cookie
	c.Redirect(http.StatusSeeOther, "/login")
}

func DeleteUser(c *gin.Context) {
	var user Users
	name := c.Param("name")
	db := database.InitDB()
	db.Where("username=?", name).Delete(&user)
	c.Redirect(303, "/wadmin")

}

func UpdateUser(c *gin.Context) {

	updateData := c.Request.FormValue("updatedata")
	var user Users
	name := c.Param("name")
	db := database.InitDB()
	db.Model(&user).Where("username=?", name).Update("username", updateData)
	c.Redirect(303, "/wadmin")
}

func CreateUser(c *gin.Context) {
	var user []Users
	var status bool = true

	FusernameN := c.Request.FormValue("username")
	Fpassword := HashPassword(c.Request.FormValue("password"))

	//database things
	db := database.InitDB()
	db.AutoMigrate(&Users{})
	db.Find(&user)

	for _, i := range user {
		if i.Username == FusernameN {
			status = false
			break
		}
	}

	if !status {
		log.Println("hello Admin , The username is already in Use")
		c.Redirect(303, "/wadmin")
		return

	}

	db.Create(&Users{Username: FusernameN, Password: Fpassword})
	log.Println("Hey Admin, Account is successfully created.")
	c.Redirect(http.StatusSeeOther, "/wadmin")

}

func IndexHandler(c *gin.Context) {
	session, _ := Store.Get(c.Request, "jwt_token")
	_, ok := session.Values["token"]
	if !ok {
		c.Redirect(303, "/login")
		return
	}
	c.Redirect(303, "/home")
}

func AdminLoged(c *gin.Context) bool {
	session, _ := Store.Get(c.Request, "admin_jwt_token")
	token, ok := session.Values["token"]
	fmt.Println(token)
	if !ok {
		return ok
	}
	return true
}

func UserLoged(c *gin.Context) bool {

	session, _ := Store.Get(c.Request, "jwt_token")
	token, ok := session.Values["token"]
	fmt.Println(token)
	if !ok {

		return ok
	}
	return true

}

func LogoutAdmin(c *gin.Context) {

	cookie, err := c.Request.Cookie("admin_jwt_token")
	if err != nil {
		c.Redirect(303, "/admin")
	}
	c.SetCookie("admin_jwt_token", "", -1, "/", "localhost", false, false)
	_ = cookie
	c.Redirect(http.StatusSeeOther, "/admin")
}
