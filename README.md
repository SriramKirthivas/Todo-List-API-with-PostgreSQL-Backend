# Todo List API with PostgreSQL Backend

This Go program implements a RESTful API for managing todo tasks, backed by a PostgreSQL database. It provides endpoints for creating, reading, updating, and deleting todo tasks through HTTP requests.

## Features

- **Create**: Users can create new todo tasks by sending a POST request to the `/create` endpoint with JSON data containing task details.
- **Read**: Users can retrieve a list of all todo tasks or a specific task by sending GET requests to the `/list` or `/show/{id}` endpoints respectively.
- **Update**: Users can update an existing todo task by sending a PUT request to the `/update/{id}` endpoint with JSON data containing the updated task details.
- **Delete**: Users can delete a todo task by sending a DELETE request to the `/delete/{id}` endpoint.

## Technologies Used

- **Gorilla Mux**: A powerful HTTP router and dispatcher for matching incoming requests to their respective handler functions.
- **pgx/pgxpool**: PostgreSQL database driver and connection pool for efficient database interactions.
- **godotenv**: A tool for loading environment variables from a `.env` file, enabling easy configuration of database connection parameters.

## Usage

1. Install Go and PostgreSQL on your system if you haven't already.
2. Clone this repository and navigate to the project directory.
3. Create a `.env` file in the project directory and set the following environment variables:
    USER=<database_user>
    PASSWORD=<database_password>
    HOST=<database_host>
    PORT=<database_port>
    DBNAME=<database_name>
4. Run the following command to build and start the server on port 8080. You can access the API endpoints using preferred HTTP Client:

```bash
go run main.go
