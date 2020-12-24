package main

import (
	"fmt"
	"net/http"
	"net/smtp"
)

//say hello ...
// func SayHello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "hello")
// }

// send email .....
func SendEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// if r.Method != "POST" {
	// 	http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
	// 	return
	// }
	fmt.Println("user sign in endpoint got hit!")
	err := r.ParseForm()
	checkErr(err)
	Username := r.FormValue("name")
	Useremail := r.FormValue("email")
	Usersubject := r.FormValue("subject")
	Usermessage := r.FormValue("message")
	fmt.Println(Username + Useremail + Usersubject + Usermessage)

	from := "fath.ibnalhaytham@gmail.com"
	password := "azerty7410."

	// Receiver email address.
	to := []string{
		"fath5yousfi@gmail.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	msg := "this message was sent from" + Username + "with email:" + Useremail + "message: " + Usermessage
	fmt.Print(msg)
	message := []byte("msg not workinggggggggggggggggggggggggggggggg!!!!!!!")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")

}
