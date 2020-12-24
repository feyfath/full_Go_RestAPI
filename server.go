package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

//checkErr
func checkErr(err error) {
	if err != nil {
		return
	}
}
func emailExists(email string) bool {
	row := db.QueryRow("select user_email from users where user_email= ?", email)
	checkErr(err)
	var temp = ""
	err := row.Scan(&temp)
	checkErr(err)
	if temp != "" {
		return true
	}
	return false
}
func getHashPassword(email string) string {
	result := ""
	row := db.QueryRow("select user_password from users where user_email= ?", email)
	_ = row.Scan(&result)
	return result
}
func getUserByID(email string) uint64 {
	result := 0
	row := db.QueryRow("select id from users where user_email=?", email)
	err := row.Scan(&result)
	checkErr(err)
	return uint64(result)
}
func createToken(userID uint64) (string, error) {
	//Creating Access Token
	_ = os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(time.Minute * 3).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	checkErr(err)
	return token, nil
}
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	_, _ = fmt.Fprint(w, "Home page!")
}
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("user sign in endpoint got hit!")
	err := r.ParseForm()
	checkErr(err)
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Println("\n the parsed form is:", r.Form)
	if emailExists(email) {
		fmt.Println("user exists!")
		err := bcrypt.CompareHashAndPassword([]byte(getHashPassword(email)), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("password is not correct")
		} else {
			fmt.Println("password is correct \n token & email were sent!")
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(`{"email":"` + email + `", `))

			w.Header().Set("Content-Type", "text/html")
			token, err := createToken(getUserByID(email))
			checkErr(err)
			w.Header().Set("Content-Type", "application/jwt")
			w.Write([]byte(`"token":"` + token + `"}`))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("invalid email address")
	}
}

func login2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("user sign in endpoint got hit!")
	err := r.ParseForm()
	checkErr(err)
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Println("\n the parsed form is:", r.Form)
	if emailExists(email) {
		fmt.Println("user exists!")
		err := bcrypt.CompareHashAndPassword([]byte(getHashPassword(email)), []byte(password))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("password is not correct")
		} else {
			w.Header().Set("Content-Type", "application/jwt")
			w.WriteHeader(http.StatusAccepted)
			fmt.Println("password is correct")
			token, err := createToken(getUserByID(email))
			checkErr(err)

			cookie := http.Cookie{
				Name:     "cookie-name",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("invalid email address")
	}
}

// func readCookie(w http.ResponseWriter, r *http.Request) {
// 	cookie, err := r.Cookie("cookie-name")
// 	if err != nil {
// 		fmt.Fprintf(w, "no cookie")
// 	}
// 	fmt.Fprintf(w, cookie.Value)
// 	fmt.Printf(cookie.Value)
// 	fmt.Fprintf(w, cookie.Name)
// 	fmt.Printf(cookie.Value)
// }

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	checkErr(err)
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	accountType := r.FormValue("account_type")
	fmt.Println("\n the form is:", r.Form)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	_, _ = db.Exec("insert into users values (default, ?, ?, ?, ?)", name, email, passwordHash, accountType)

	fmt.Println("user registered with success \n token & email were sent!")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"email":"` + email + `", `))

	w.Header().Set("Content-Type", "text/html")
	token, err := createToken(getUserByID(email))
	checkErr(err)
	w.Header().Set("Content-Type", "application/jwt")
	w.Write([]byte(`"token":"` + token + `"}`))

}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
	id, ok := r.URL.Query()["id"]
	if !ok {
		log.Println("Url Param 'key' is missing")
		return
	}
	fmt.Print("delete user with id:", id)
	_, _ = db.Exec("delete from users where id = ?", id)
}

func handler() {
	fmt.Printf("hello handler!")
	db, err = sql.Open("mysql", "root:8520@tcp(127.0.0.1:3306)/world")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("successfully connected to database.")
	}
	defer db.Close()
	http.HandleFunc("/", home)
	http.HandleFunc("/SendEmail", SendEmail)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login2", login2)
	http.HandleFunc("/login", login)
	http.HandleFunc("/user/:id", deleteUser)
	// http.HandleFunc("/GetAvatar", GetAvatar)
	http.HandleFunc("/SetAvatar", SetAvatar)
	// http.HandleFunc("/rc", readCookie)
	// http.Handle("/", http.FileServer(http.Dir("dist")))
	log.Fatal(http.ListenAndServe(":8082", nil))

}
func main() {
	godotenv.Load(".env")
	secretKey := os.Getenv("ACCESS_SECRET")
	fmt.Print("ACCESS_SECRET is:" + secretKey + "\n")
	handler()
}
