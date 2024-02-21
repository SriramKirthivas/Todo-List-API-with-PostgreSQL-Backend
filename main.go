package main

// Context - defines cancellation signals or deadlines across various API Calls
// mux - HTTP request multiplexer- implements request router and dispatcher
//pgxpool - connection pool for pgx
//godotenv - use to load environment variable in this case
import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	//"strconv"

	"github.com/gorilla/mux"
	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

// Establishing a variable to be used for sql queries
var db *pgxpool.Pool

// TODO structure with variables ID,Title,Content.
type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done string `json:"done"`
}

// type ToDoWeb struct {
// 	ToDoData []Todo
// }

// This function is used to initialize the variables required to access the database

func initDB() {
	// load environment variable. It returns error value and it is assigned to err. If error is there, it prints the out error inside log.Fatal
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variable file")
	}
	// We get the environment variables stored in .env via os.Getenv and associate to a variable
	user := os.Getenv("USER")
	passwd := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	db_name := os.Getenv("DBNAME")
	// This is the database URL. Sprintf is used for storing the output of printf to a variable. Here the conn stored the database URL.
	conn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, passwd, host, port, db_name)
	// We connect to the database using pgxpool.Connect with parameters context.Background() and the URL. The returned values are the connection and error.
	connection, errDataBase := pgxpool.Connect(context.Background(), conn)
	if errDataBase != nil {
		log.Fatal("Error connecting to the database ", errDataBase)
	}
	// We assign the connection value to db
	db = connection
}

// Main function, is used to create a router and connect each functionality to a paticular path with each methods to do.
func main() {
	// The initialization function is called
	initDB()
	// We create a new Router instance and assign it to r
	r := mux.NewRouter()
	//r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("Downloads/http_server"))))
	// For different path, we assign different functionalities and methods as well
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/create", createToDo).Methods("POST")
	r.HandleFunc("/update/{id}", updateToDo).Methods("PUT")
	r.HandleFunc("/delete/{id}", deleteToDo).Methods("DELETE")
	r.HandleFunc("/list", listToDo).Methods("GET")
	r.HandleFunc("/show/{id}", showToDo).Methods("GET")
	log.Printf("Server running on port 8080")
	// The server is running on port 8080 with a paticular server max in this case r is used.
	log.Fatal(http.ListenAndServe(":8080", r))

}

// CreateToDo() is used to create a todo list and add it to the database
func createToDo(w http.ResponseWriter, r *http.Request) {
	// A variable todo of type Todo struct is initialized
	var todo Todo
	// NewDecoder is used to parse the json request body and decode function, decodes the json request and point it to todo. The return value is an error.
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
	// QueryRow function, establishes a connection and executes a query that returns atmost 1 row. It returns a row.
	row := db.QueryRow(r.Context(), "INSERT INTO todo(id,task,done) VALUES ($1,$2,$3) RETURNING id", todo.ID, todo.Task, todo.Done)
	// Scan works like Rows but, when no rows are returned it returns ErrNoRows error.
	if err := row.Scan(&todo.ID); err != nil {
		http.Error(w, "Error creating todo"+err.Error(), http.StatusInternalServerError)
		return
	}
	// Header() returns the header map that is sent by the writeheader. And we set the elements associated with keys as a single value. Here, we set Content-Type as application/json
	w.Header().Set("Content-Type", "application/json")
	// NewEncoder returns a new encoder that writes to w. Encode writes json encoding of todo to stream
	json.NewEncoder(w).Encode(todo)
}

// UpdateToDo i used to update a row based on id value
func updateToDo(w http.ResponseWriter, r *http.Request) {
	// Vars returns the route variable of a paticular request. Like in update/{id} if we use update/1, it returns out the entire map of id 1 and assign it to vars.
	vars := mux.Vars(r)
	//Atoi function converts a ascii value to integer. It returns an integer and error.
	todoid, _ := strconv.Atoi(vars["id"])
	// Initialize updateToDo variable of type Todo struct
	var updateTodo Todo
	// NewDecoder is used to parse the json request body and decode function, decodes the json request and point it to updateTodo. The return value is an error.
	if err := json.NewDecoder(r.Body).Decode(&updateTodo); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
	}
	// Exec acquires connection with the database and executes the given instruction. It returns a command tag and error. Here we use UPDATE instruction to update the title and content based on id
	_, err := db.Exec(r.Context(), "UPDATE todo SET task=$1,done=$2 WHERE id=$3", updateTodo.Task, updateTodo.Done, todoid)
	if err != nil {
		http.Error(w, "Error updating todo "+err.Error(), http.StatusInternalServerError)
		return
	}
	// After updating, we send the appropriate status code to the response header
	w.WriteHeader(http.StatusOK)
}
func deleteToDo(w http.ResponseWriter, r *http.Request) {
	// Vars returns the route variable of a paticular request. Like in update/{id} if we use update/1, it returns out the entire map of id 1 and assign it to vars.
	vars := mux.Vars(r)
	//Atoi function converts a ascii value to integer. It returns an integer and error.
	todoId, _ := strconv.Atoi(vars["id"])
	// Exec acquires connection with the database and executes the given instruction. It returns a command tag and error. Here, we delete a paticular row based on id.
	_, err := db.Exec(r.Context(), "DELETE FROM todo WHERE id = $1", todoId)
	if err != nil {
		http.Error(w, "Error executing todo "+err.Error(), http.StatusInternalServerError)
		return
	}
	// After updating, we send the appropriate status code to the response header
	w.WriteHeader(http.StatusOK)
}
func listToDo(w http.ResponseWriter, r *http.Request) {
	// Query function executes the sql and returns the row and error
	row, err := db.Query(r.Context(), "SELECT id,task,done FROM todo ORDER BY id ASC")
	if err != nil {
		http.Error(w, "Error fetching the todo list", http.StatusInternalServerError)
		return
	}
	// Close - closes the rows so that connection is available to use again
	defer row.Close()
	// Initialize todos variable of a silce of struct Todo. It is initialized to lise each ToDo list
	var todos []Todo
	// Next function prepares the next row for reading.
	for row.Next() {
		// Initialize todo of Todo struct. Scan reads the value of a paticular row positionally. It returns an error
		var todo Todo
		if err := row.Scan(&todo.ID, &todo.Task, &todo.Done); err != nil {
			http.Error(w, "Error scanning todo", http.StatusInternalServerError)
			return
		}
		//Append the value of todo to todos
		todos = append(todos, todo)
	}
	// Header() returns the header map that is sent by the writeheader. And we set the elements associated with keys as a single value. Here, we set Content-Type as application/json
	w.Header().Set("Content-Type", "application/json")
	// NewEncoder returns a new encoder that writes to w. Encode writes json encoding of todos to stream
	json.NewEncoder(w).Encode(todos)
}
func showToDo(w http.ResponseWriter, r *http.Request) {
	// Vars returns the route variable of a paticular request. Like in update/{id} if we use update/1, it returns out the entire map of id 1 and assign it to vars.
	vars := mux.Vars(r)
	//Atoi function converts a ascii value to integer. It returns an integer and error.
	toDoid, _ := strconv.Atoi(vars["id"])
	//Initialize todo variable of type Todo struct
	var todo Todo
	// QueryRow function, establishes a connection and executes a query that returns atmost 1 row. It returns a row. And we scan it with id,title and content. It returns a error
	err := db.QueryRow(r.Context(), "SELECT id,task,done FROM todo WHERE id=$1", toDoid).Scan(&todo.ID, &todo.Task, &todo.Done)
	if err == sql.ErrNoRows {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching the todo list", http.StatusInternalServerError)
		return
	}
	// Header() returns the header map that is sent by the writeheader. And we set the elements associated with keys as a single value. Here, we set Content-Type as application/json
	w.Header().Set("Content-Type", "application/json")
	// NewEncoder returns a new encoder that writes to w. Encode writes json encoding of todo to stream
	json.NewEncoder(w).Encode(todo)
}
