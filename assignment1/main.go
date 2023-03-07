package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	auth "github.com/Bektemis/golang_ass_1/authorization"
	search "github.com/Bektemis/golang_ass_1/item_search"
	pck "github.com/Bektemis/golang_ass_1/pck"
	rate "github.com/Bektemis/golang_ass_1/ratings"
	regist "github.com/Bektemis/golang_ass_1/registration"
	jwt "github.com/dgrijalva/jwt-go"
)

type database struct {
	items *pck.DatabaseItems
	users *pck.DatabaseUsers
}

var db database = database{items: &pck.DatabaseItems{Items: []pck.Item{}}, users: &pck.DatabaseUsers{Users: []pck.User{}}}

// used to encrypt/decrypt JWT tokens.
var jwtTokenSecret = "secretjwt"

const dashBoardPage = `<html><body>
{{if .Username}}
<p><b>{{.Username}}</b>, welcome to your dashboard! <a href="/logout">Logout!</a></p>
{{else}}
<p>Either your JSON Web token has expired or you've logged out! <a href="/login">Login</a></p>
{{end}}
<form>
<input type="text" name="search" placeholder="Search">
<input type="submit" name ="Search" value="Search">
<form><br>
<form>
<input type="text" name="item" placeholder="Pick item">
<input type="text" name="rating" placeholder="Rate item">
<input type="submit" name ="Rate" value="Rate">
<form><br>
<form>
    <input type="submit" name = "by-rating" value="Sort By Ratings"><br>
    <input type="submit" name = "by-price" value="Sort By Price">
</form><br>
{{.list}}
	
</body></html>`

const logUserPage = `<html><body>
  {{if .LoginError}}<p style="color:red">Either username or password is not in our record! Sign Up?</p>{{end}}
  <form method="post" action="/login">
  {{if .Username}}
  <p><b>{{.Username}}</b>, you're already logged in! <a href="/logout">Logout!</a></p>
  {{else}}
  <label>Username:</label>
  <input type="text" name="Username"><br>
  <label>Password:</label>
  <input type="password" name="Password">
  <input type="submit" name="Login" value="Let me in!">
  {{end}}
  </form>
  <br>
  <form action="/regist">
    <input type="submit" value="Registration" />
  </form>
  </body></html>`

const RegistUserPage = `<html><body>
  {{if .RegistError}}<p style="color:red">Username is not available</p>{{end}}
  <form method="post" action="/regist">
  <label>Username:</label>
  <input type="text" name="Username"><br>
  <label>Password:</label>
  <input type="password" name="Password"><br>
  <input type="submit" name="" value="Register">
  </form>
  <br>
  </body></html>`

func DashBoardPageHandler(w http.ResponseWriter, r *http.Request) {
	conditionsMap := map[string]interface{}{}

	// THIS PAGE should ONLY be accessible to those users that logged in

	// check if user already logged in
	username, _ := ExtractTokenUsername(r)

	if username != "" {
		conditionsMap["Username"] = username
	}
	conditionsMap["list"] = db.items.GetListOfItems()
	if r.FormValue("by-rating") != "" {
		db.items.FilterByRatings(true)
		conditionsMap["list"] = db.items.GetListOfItems()
	}
	if r.FormValue("by-price") != "" {
		db.items.FilterByPrice(true)
		conditionsMap["list"] = db.items.GetListOfItems()
	}
	if r.FormValue("search") != "" && r.FormValue("Search") != "" {
		searchValue := r.FormValue("search")
		conditionsMap["list"] = search.ItemSearch(searchValue, db.items)
	}
	if r.FormValue("rating") != "" && r.FormValue("Rate") != "" && r.FormValue("item") != "" {
		rateValue := r.FormValue("rating")
		item := r.FormValue("item")
		x, err := strconv.Atoi(rateValue)
		if err != nil {
			fmt.Println("Error during conversion")
			return
		}
		rate.GiveRating(x, item, db.items)
		for _, i := range db.items.GetListOfItems() {
			fmt.Println(i)
		}
		conditionsMap["list"] = db.items.GetListOfItems()
	}

	if err := dashboardTemplate.Execute(w, conditionsMap); err != nil {
		log.Println(err)
	}
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {

	conditionsMap := map[string]interface{}{}

	// verify username and password
	if r.FormValue("Login") != "" && r.FormValue("Username") != "" {
		username := r.FormValue("Username")
		password := r.FormValue("Password")

		if !auth.SignIn(username, password, db.users) {
			log.Println("Either username or password is wrong")
			conditionsMap["LoginError"] = true
		} else {
			log.Println("Logged in :", username)
			conditionsMap["Username"] = username
			conditionsMap["LoginError"] = false

			// create a new JSON Web Token and redirect to dashboard
			tokenString, err := createToken(username)

			if err != nil {
				log.Println(err) // of course, this is too simple, your program should prevent login if token cannot be generated!!
				os.Exit(1)
			}

			// create the cookie for client(browser)
			expirationTime := time.Now().Add(1 * time.Hour) // cookie expired after 1 hour

			cookie := &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			}

			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/dashboard", http.StatusFound)
		}

	}

	if err := logUserTemplate.Execute(w, conditionsMap); err != nil {
		log.Println(err)
	}
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {

	conditionsMap := map[string]interface{}{}

	// check if user already logged in
	username, _ := ExtractTokenUsername(r)

	if username != "" { // user already logged in!
		conditionsMap["Username"] = username
		conditionsMap["LoginError"] = false
	}

	// verify username and password
	if r.FormValue("Password") != "" && r.FormValue("Username") != "" {
		username := r.FormValue("Username")
		password := r.FormValue("Password")

		if !regist.Register(username, password, db.users) {
			log.Println("Username is not available")
			conditionsMap["LoginError"] = true
		} else {
			log.Println("Registered:", username, password)
			conditionsMap["Username"] = username
			conditionsMap["LoginError"] = false

			// create a new JSON Web Token and redirect to dashboard
			tokenString, err := createToken(username)

			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			// create the cookie for client(browser)
			expirationTime := time.Now().Add(1 * time.Hour) // cookie expired after 1 hour

			cookie := &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			}

			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/dashboard", http.StatusFound)
		}

	}

	if err := registerUserTemplate.Execute(w, conditionsMap); err != nil {
		log.Println(err)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(w, &c)

	w.Write([]byte("Old cookie deleted. Logged out!\n"))
}

func createToken(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username                            //embed username inside the token string
	claims["expired"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtTokenSecret))
}

func ExtractTokenUsername(r *http.Request) (string, error) {

	// get our token string from Cookie
	biscuit, err := r.Cookie("token")

	var tokenString string
	if err != nil {
		tokenString = ""
	} else {
		tokenString = biscuit.Value
	}

	// abort
	if tokenString == "" {
		return "", nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtTokenSecret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username := fmt.Sprintf("%s", claims["username"]) // convert to string
		if err != nil {
			return "", err
		}
		return username, nil
	}
	return "", nil
}

var dashboardTemplate = template.Must(template.New("").Parse(dashBoardPage))
var logUserTemplate = template.Must(template.New("").Parse(logUserPage))
var registerUserTemplate = template.Must(template.New("").Parse(RegistUserPage))

func main() {
	db.items.Items = append(db.items.Items, pck.Item{Name: "item1", Price: 100, Rating: 5, HaveRated: 1})
	db.items.Items = append(db.items.Items, pck.Item{Name: "item2", Price: 200, Rating: 4, HaveRated: 1})
	db.items.Items = append(db.items.Items, pck.Item{Name: "item3", Price: 100, Rating: 3, HaveRated: 1})
	db.items.Items = append(db.items.Items, pck.Item{Name: "item4", Price: 600, Rating: 1, HaveRated: 1})
	regist.Register("Almat", "Bektemis", db.users)
	fmt.Println("Server starting, point your browser to localhost:8080/login to start")
	http.HandleFunc("/login", LoginPageHandler)
	http.HandleFunc("/dashboard", DashBoardPageHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/regist", RegisterPageHandler)
	http.ListenAndServe(":8080", nil)
}
