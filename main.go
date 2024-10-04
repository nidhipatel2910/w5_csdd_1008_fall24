package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)

// Task struct definition
type Task struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"` // "pending" or "completed"
}

var tasks []Task
var nextID = 1

// Create a new task
func createTask(w http.ResponseWriter, r *http.Request) {
    var task Task
    json.NewDecoder(r.Body).Decode(&task)
    task.ID = nextID
    nextID++
    tasks = append(tasks, task)
    json.NewEncoder(w).Encode(task)
}

// Get all tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(tasks)
}

// Get a task by ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])
    for _, task := range tasks {
        if task.ID == id {
            json.NewEncoder(w).Encode(task)
            return
        }
    }
    http.Error(w, "Task not found", http.StatusNotFound)
}

// Update a task by ID
func updateTask(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])
    for index, task := range tasks {
        if task.ID == id {
            json.NewDecoder(r.Body).Decode(&task)
            task.ID = id // Maintain the same ID
            tasks[index] = task
            json.NewEncoder(w).Encode(task)
            return
        }
    }
    http.Error(w, "Task not found", http.StatusNotFound)
}

// Delete a task by ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])
    for index, task := range tasks {
        if task.ID == id {
            tasks = append(tasks[:index], tasks[index+1:]...)
            json.NewEncoder(w).Encode(tasks)
            return
        }
    }
    http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/tasks", createTask).Methods("POST")
    r.HandleFunc("/tasks", getTasks).Methods("GET")
    r.HandleFunc("/tasks/{id}", getTaskByID).Methods("GET")
    r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
    r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

    fmt.Println("Server running on port 8080")
    http.ListenAndServe(":8080", r)
}
