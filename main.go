package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Task represents a simple task with an ID, Title, Description, and Status
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // Task status: "pending" or "completed"
}

var (
	tasks    = make(map[int]Task) // In-memory storage for tasks
	nextID   = 1                  // Keeps track of the next task ID
	taskLock sync.Mutex           // Ensures safe access to tasks in concurrent operations
)

// Create a new task (POST /tasks)
func createTask(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the task from the request body
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lock before modifying shared data
	taskLock.Lock()
	defer taskLock.Unlock()

	// Assign a unique ID and save the task
	newTask.ID = nextID
	nextID++
	tasks[newTask.ID] = newTask

	// Respond with the created task in JSON format
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// Get a task by ID (GET /tasks/{id})
func getTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr) // Convert the string ID to an integer
	if err != nil || id < 1 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Lock to safely access the task
	taskLock.Lock()
	defer taskLock.Unlock()

	// Find the task by ID
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Respond with the task in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Update a task by ID (PUT /tasks/{id})
func updateTask(w http.ResponseWriter, r *http.Request) {
	// Check if the method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from the URL
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Parse the updated task from the request body
	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lock before modifying the task
	taskLock.Lock()
	defer taskLock.Unlock()

	// Check if the task exists
	_, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Update the task and respond with the updated task in JSON format
	updatedTask.ID = id
	tasks[id] = updatedTask
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

// Delete a task by ID (DELETE /tasks/{id})
func deleteTask(w http.ResponseWriter, r *http.Request) {
	// Check if the method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract task ID from the URL
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Lock before modifying tasks
	taskLock.Lock()
	defer taskLock.Unlock()

	// Check if the task exists
	_, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Delete the task
	delete(tasks, id)

	// Respond with no content status
	w.WriteHeader(http.StatusNoContent)
}

// List all tasks (GET /tasks)
func listTasks(w http.ResponseWriter, r *http.Request) {
	// Lock to safely access the tasks
	taskLock.Lock()
	defer taskLock.Unlock()

	// Collect all tasks in a slice
	var taskList []Task
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	// Respond with the task list in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}

// Main function to set up routes and start the server
func main() {
	// Route for listing and creating tasks
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			listTasks(w, r)
		} else if r.Method == http.MethodPost {
			createTask(w, r)
		} else {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	})

	// Route for accessing, updating, and deleting a task by ID
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTask(w, r)
		} else if r.Method == http.MethodPut {
			updateTask(w, r)
		} else if r.Method == http.MethodDelete {
			deleteTask(w, r)
		} else {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	})

	// Start the server on port 8090
	fmt.Println("Server is running on port 8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
