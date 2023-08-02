# Player Database API.

This is a simple web application that provides an API to manage player records in a PostgreSQL database. The application is built using the Go programming language and utilizes the Gin web framework to handle HTTP requests.

## Installation and Setup.

 1. Make sure you have Go (Golang) installed on your machine.

 2. Clone this repository to your local machine:

     `git clone https://github.com/your-username/player-database-api.git`

 3. Install the required dependencies:

     `go mod download`

 4. Set up the PostgreSQL database:

    * Install PostgreSQL on your machine if you haven't already.
    * Create a new database and configure the connection details in the code (`DB_DSN` constant in  `main.go`).   

## Docker compose command

```bash

docker-compose -f docker-compose-postgres.yml -f docker-compose-app.yml up -d

```

## Usage.

To run the application, execute the following command from the root directory of the project:

     `go run main.go`

The application will start and listen on port 8080 for incoming HTTP requests.

## API Endpoints.

The application provides the following API endpoints:

 1. Create Player: Add a new player record to the database.


     `POST /players`

 2. Update Player: Update an existing player record in the database.

     `PUT /players/:id`

 3. Delete Player: Remove a player record from the database.

     `DELETE /players/:id`

 4. List Players: Get a list of all players in descending order of their IDs.

     `GET /players`

 5. Get Player by Rank: Get a player based on their rank (ranked by their score in descending order).

     `GET /players/rank/:val`

 6. Get Random Player: Get a random player from the database.

     `GET /players/random`

## Request and Response Format.

   * The API expects JSON payload for creating and updating players.

   * Successful responses will be returned in JSON format with appropriate HTTP status codes.

   * Error responses will contain an `error` field with an error message and the corresponding HTTP status code.

## Example Usage.

Assuming the application is running locally on port 8080, you can use tools like `curl` or Postman to interact with the API.

 1. To create a new player:

     `curl -X POST -H "Content-Type: application/json" -d '{"name": "Nitish", "country": "IN", "score": 100}' http://localhost:8080/players`

 2. To update an existing player with ID 1:

     `curl -X PUT -H "Content-Type: application/json" -d '{"name": "Nitish_Kumar", "score": 150}' http://localhost:8080/players/1`

 3. To delete a player with ID 1:

     `curl -X DELETE http://localhost:8080/players/1`

 4. To list all players:

     `curl http://localhost:8080/players`

 5. To get a player by rank (e.g., rank 3):

     `curl http://localhost:8080/players/rank/3`

 6. To get a random player:

     `curl http://localhost:8080/players/random`

## Important Notes.

   * Ensure that you have set up the PostgreSQL database and provided the correct database connection details in the code (`DB_DSN` constant).

   * This application is for demonstration purposes and may lack certain production-ready features such as authentication and proper error handling.

   * Feel free to modify the code and add additional features as needed for your use case.

