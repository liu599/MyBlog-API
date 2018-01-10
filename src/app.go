package main
// http://blog.csdn.net/pmlpml/article/details/78539261
import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
	"log"
	"net/http"
	"encoding/json"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"github.com/rs/cors"
	// "reflect"
	"github.com/codegangsta/negroni"
	"reflect"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Handler http.Handler
}

func (a *App) Run (addr string) {
	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(ValidateTokenMiddleware))
	n.UseHandler(a.Handler)
	n.Run(":12450")
	// log.Fatal(http.ListenAndServe(":12450", a.Handler))
}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/token.get", a.GenerateTokenMiddleware).Methods("POST")

	a.Router.HandleFunc("/user/{uid:[0-9]+}", a.getUserDetail).Methods("POST")

	a.Router.HandleFunc("/auth/users", a.getUsers).Methods("POST")
	a.Router.HandleFunc("/auth/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/auth/user", a.getUser).Methods("POST")
	a.Router.HandleFunc("/auth/user", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/auth/user", a.deleteUser).Methods("DELETE")

	a.Router.HandleFunc("/categories", a.getCategories).Methods("POST")
	a.Router.HandleFunc("/auth/category", a.createCategory).Methods("POST")
	a.Router.HandleFunc("/category/{cid:[0-9]+}", a.getCategory).Methods("POST")
	a.Router.HandleFunc("/auth/category", a.updateCategory).Methods("PUT")
	a.Router.HandleFunc("/auth/category", a.deleteCategory).Methods("DELETE")

	a.Router.HandleFunc("/posts", a.getPosts).Methods("POST")
	a.Router.HandleFunc("/post", a.getPost).Methods("POST")
	a.Router.HandleFunc("/posts/category/{cid:[0-9a-zA-Z]+}", a.getPostsByCategory).Methods("POST")
	a.Router.HandleFunc("/auth/post", a.deletePost).Methods("DELETE")
	a.Router.HandleFunc("/auth/post", a.updatePost).Methods("PUT")
	a.Router.HandleFunc("/auth/post", a.createPost).Methods("POST")

	a.Router.HandleFunc("/comments/{pid:[0-9]+}", a.getComments).Methods("POST")
	a.Router.HandleFunc("/comment/{coid:[0-9]+}", a.getComment).Methods("POST")
	a.Router.HandleFunc("/auth/comment/{coid:[0-9]+}", a.deleteComment).Methods("DELETE")
	a.Router.HandleFunc("/auth/comment/{coid:[0-9]+}", a.updateComment).Methods("PUT")
	a.Router.HandleFunc("/auth/comment", a.createComment).Methods("POST")

	a.Router.HandleFunc("/pxxxx/{pid:[0-9]+}", a.postWriter).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"POST"},
		AllowCredentials: true,
	})

	handler := c.Handler(a.Router)

	a.Handler = handler

	// fmt.Println("type:", reflect.TypeOf(handler))
}


func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	// JUST ADD THE ROUTES
	a.initializeRoutes()
}


func validateUserInfo(r *http.Request) string {
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			// fmt.Printf("%v: %v \n", name, h)
			if name == "authorization" && h == os.Getenv("NEKOHAND_AUTHORIZATION") {
				return "right"
			}
			// request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	return "not right"
}


func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// fmt.Printf("%v\n", payload)
	// response, _ := json.Marshal(payload)
	fmt.Println("type:", reflect.TypeOf(payload))
	body := make(map[string]interface{})
	body["data"] = payload
	if code == 200 || code == 201 {
		body["msg"] = "OK"
		body["code"] = "0"
		body["success"] = true
	} else {
		body["msg"] = "ERROR"
		body["code"] = "1"
		body["success"] = false
	}

	response, _ := json.Marshal(body)

	// fmt.Printf(string(response))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
	w.WriteHeader(code)
	w.Write(response)
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}


// 生成32位MD5
func MD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 类型函数
func type_cast(general interface{}) {
	switch general.(type) {
	case int :
		fmt.Println("is int")
	case float64 :
		fmt.Println("is float64")
	case string :
		fmt.Println("is string")
	default :
		fmt.Println("Unknown Type")
	}
}