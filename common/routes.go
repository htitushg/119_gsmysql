package common

import (
	"net/http"
)

func Routes() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/Login", Login)
	http.HandleFunc("/Apropos", Apropos)
	http.HandleFunc("/Signin", Signin)
	http.HandleFunc("/Refresh", Refresh)
	http.HandleFunc("/Logout", Logout)
	http.HandleFunc("/AfficheUserInfo", AfficheUserInfo)
	http.HandleFunc("/Register", Register)
	http.HandleFunc("/RegisterPost", RegisterPost)
}
