package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Player struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Score   int    `json:"score"`
}

var DB *sql.DB

func main() {
	createDBConnection()
	defer DB.Close()
	r := gin.Default()
	setupRouters(r)
	r.Run()

}

func setupRouters(r *gin.Engine) {

	r.POST("/players", createPlayer)
	r.PUT("/players/:id", updatePlayer)
	r.DELETE("/players/:id", deletePlayer)
	r.GET("/players", listPlayers)
	r.GET("/players/rank/:val", getPlayerByRank)
	r.GET("/players/random", getRandomPlayer)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}

func createPlayer(c *gin.Context) {
	var player Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if player.Name == "" || len(player.Name) > 15 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}

	if len(player.Country) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country code"})
		return
	}

	// Insert the new player into the database
	row := DB.QueryRow("INSERT INTO players (name, country, score) VALUES ($1, $2, $3) RETURNING id",
		player.Name, strings.ToUpper(player.Country), player.Score)

	if err := row.Scan(&player.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player created successfully", "player": player})
}

func updatePlayer(c *gin.Context) {
	idStr := c.Param("id")

	fmt.Println(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	var player Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	fmt.Println("printing the struct")
	fmt.Printf("%+v\n", player)

	var count int

	// Check if the player exists in the database
	row := DB.QueryRow("SELECT COUNT (*) FROM players WHERE id = $1", id)
	if err := row.Scan(&count); err != nil {
		if count != 0 {
			c.JSON(http.StatusFound, gin.H{"error": "id do not exists"})
		}
		return
	}

	// Only update name and score
	if player.Name != "" {
		if len(player.Name) > 15 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
			return
		}
	}
	if player.Score != 0 {
		_, err = DB.Exec("UPDATE players SET name = $1, score = $2 WHERE id = $3", player.Name, player.Score, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update player"})
			return
		}
	}

	rows, err := DB.Query("SELECT * FROM players where id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&player.ID, &player.Name, &player.Country, &player.Score)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"message": "Player updated successfully", "player": player})
}

func deletePlayer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID"})
		return
	}

	// Check if the player exists in the database
	row := DB.QueryRow("SELECT id FROM players WHERE id = $1", id)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch player data"})
		}
		return
	}

	_, err = DB.Exec("DELETE FROM players WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})
}

func listPlayers(c *gin.Context) {
	rows, err := DB.Query("SELECT * FROM players ORDER BY id DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}
	defer rows.Close()

	playerList := make([]Player, 0)
	for rows.Next() {
		var player Player
		err := rows.Scan(&player.ID, &player.Name, &player.Country, &player.Score)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
			return
		}
		playerList = append(playerList, player)
	}

	c.JSON(http.StatusOK, playerList)
}

func getPlayerByRank(c *gin.Context) {
	rankStr := c.Param("val")
	rank, err := strconv.Atoi(rankStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rank"})
		return
	}

	if rank < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rank"})
		return
	}

	rows, err := DB.Query("SELECT * FROM players ORDER BY score DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}
	defer rows.Close()

	var player Player
	for i := 1; i <= rank; i++ {
		if !rows.Next() {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}
		err := rows.Scan(&player.ID, &player.Name, &player.Country, &player.Score)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
			return
		}
	}

	c.JSON(http.StatusOK, player)
}

func getRandomPlayer(c *gin.Context) {
	rows, err := DB.Query("SELECT * FROM players")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}
	defer rows.Close()

	playerList := make([]Player, 0)
	for rows.Next() {
		var player Player
		err := rows.Scan(&player.ID, &player.Name, &player.Country, &player.Score)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
			return
		}
		playerList = append(playerList, player)
	}

	if len(playerList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No players found"})
		return
	}

	// Generate a random index to select a random player from the list
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(playerList))

	c.JSON(http.StatusOK, playerList[index])
}
