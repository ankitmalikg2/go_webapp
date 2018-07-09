package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var templates *template.Template
var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	muxRoute := mux.NewRouter()
	muxRoute.HandleFunc("/hello", helloHandler).Methods("GET")
	muxRoute.HandleFunc("/bye", byeHandler).Methods("GET")
	muxRoute.HandleFunc("/ankit", testHandler).Methods("GET")
	muxRoute.HandleFunc("/", indexGetHandler).Methods("GET")
	muxRoute.HandleFunc("/", indexPostHandler).Methods("POST")

	//fs := http.FileServer(http.Dir("./static/"))
	//muxRoute.PathPrefix("static").Handler(http.StripPrefix("/static/", fs))
	muxRoute.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", muxRoute)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Fprint(res, "Hello from go webapp")
}

func byeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Bye !")
}

func indexGetHandler(res http.ResponseWriter, req *http.Request) {
	comments, err := redisClient.LRange("comments", 0, -1).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(res, "index.html", comments)
}

func indexPostHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	comment := req.PostForm.Get("comment")
	redisClient.LPush("comments", comment)
	http.Redirect(res, req, "/", 302)
}

func testHandler(res http.ResponseWriter, req *http.Request) {
	templates.Execute(res, "ankit malik")
}
