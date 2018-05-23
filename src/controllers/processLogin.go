package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"

	user "github.com/ninadakolekar/aizant-dms/src/user"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// ProcessLogin ... Processes Login form
func ProcessLogin(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	redirectTarget := "/"

	if user.ValidateLoginCredentials(username, password) {

		// Check Credentials are correct or not
		// user.AuthCredentials(username,password)

		setSession(username, w)
		redirectTarget = "/dashboard"

	}

	http.Redirect(w, r, redirectTarget, 302)
}

func setSession(username string, response http.ResponseWriter) {

	value := map[string]string{
		"name": username,
	}

	encoded, err := cookieHandler.Encode("session", value)

	if err != nil {
		fmt.Println("Error setSession Line 36: ", err) // Debug
		return
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
	}

	http.SetCookie(response, cookie)
}