package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Todo struct {
	Id        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var Todos = []Todo{
	{Id: "1", Item: "Wake Up", Completed: false},
	{Id: "2", Item: "Sleep", Completed: false},
	{Id: "3", Item: "Learn Golang", Completed: false},
}

func getAllTodos(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(Todos)
	if err != nil {
		http.Error(res, "Encoding Failed", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func getTodoById(res http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/todos/")
	if id == "" {
		http.Error(res, "Invalid Id", http.StatusBadRequest)
		return
	}
	for i := range Todos {
		if Todos[i].Id == id {
			err := json.NewEncoder(res).Encode(Todos[i])
			if err != nil {
				http.Error(res, "Encoding Failed", http.StatusInternalServerError)
				return
			}
			res.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(res, "Invalid Id", http.StatusBadRequest)
	return
}

func addTodo(res http.ResponseWriter, req *http.Request) {
	var todo Todo
	err := json.NewDecoder(req.Body).Decode(&todo)
	if err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusCreated)
	Todos = append(Todos, todo)
}

func doneStatus(res http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/todos/")
	if id == "" {
		http.Error(res, "Invalid Id", http.StatusBadRequest)
		return
	}
	var val Todo
	err := json.NewDecoder(req.Body).Decode(&val)
	if err != nil {
		http.Error(res, "Invalid Json Format", http.StatusBadRequest)
		return
	}

	for i := range Todos {
		if Todos[i].Id == id {
			Todos[i].Completed = val.Completed
			return
		}
	}

	http.Error(res, "Resource Not Found", http.StatusNotFound)
	return
}

func main() {

	// Lets create the routes.
	http.HandleFunc("GET /todos", getAllTodos)
	http.HandleFunc("GET /todos/", getTodoById)
	http.HandleFunc("POST /todos", addTodo)
	http.HandleFunc("PATCH /todos/", doneStatus)

	// Lets start the server

	http.ListenAndServe("localhost:8080", nil)
}
