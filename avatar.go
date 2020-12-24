package main

import (
	"fmt"
	"net/http"
	"time"
)

// func GetAvatar(w http.ResponseWriter, r *http.Request) {
// 	row := db.QueryRow("select avatar from users where id= 2")
// }
func SetAvatar(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()
	userId := r.FormValue("id")
	avatar := r.FormValue("avatar")
	time.After(time.Second * 30)
	fmt.Println("Timed out")

	fmt.Println("userid=" + userId + "user Avatar" + avatar)
	db.Exec("update users set img_str = ? where id = ?", avatar, userId)
	time.After(time.Second * 30)
	fmt.Println("Timed out")
}
