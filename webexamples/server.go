package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

type Todo struct {
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Data       any    `json:"data"`
	Details    string `json:"details"`
}

var (
	todoItems    = make(map[int]Todo)
	todoTemplate = template.Must(template.ParseFiles("templates/todo-form.html"))
)

func logWrapper(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", logWrapper(homeHandler))
	router.HandleFunc("/todo", logWrapper(todoHandler))

	handleBookRoutes(router)
	handleAPIRoutes(router)

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":80", router)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my website!")
}

func handleBookRoutes(router *mux.Router) {
	booksRouter := router.PathPrefix("/books").Subrouter()

	booksRouter.HandleFunc("/", logWrapper(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Books Homepage")
	}))

	booksRouter.HandleFunc("/{title}/page/{page}", logWrapper(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	}))
}

func handleAPIRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/{api:api(?:\\/)?}").Subrouter()

	apiRouter.HandleFunc("", logWrapper(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Available Endpoints:\n- /todo-new [POST, JSON]\n- /todo-list [GET, JSON]")
	}))

	apiRouter.HandleFunc("/todo-new", logWrapper(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var todo Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{StatusCode: http.StatusBadRequest, Details: err.Error()})
			return
		}

		id := len(todoItems)
		todoItems[id] = todo

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{StatusCode: http.StatusCreated, Details: "", Data: todo})
	}))

	apiRouter.HandleFunc("/todo-list", logWrapper(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{StatusCode: http.StatusCreated, Details: "", Data: todoItems})
	}))
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		todoTemplate.Execute(w, nil)
		return
	}

	description := r.FormValue("description")
	done := rand.Int()%2 == 0
	todoItems[len(todoItems)] = Todo{description, done}

	todoTemplate.Execute(w, struct {
		Success   bool
		TodoItems map[int]Todo
	}{true, todoItems})
}
