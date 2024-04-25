package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	rdb *redis.Client
	ctx = context.Background()
)

func initDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=Qasaqayj7 dbname=db_host sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully!")
	return db
}

func initRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
	return client
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")
	// Check cache
	val, err := rdb.Get(ctx, id).Result()
	if err == redis.Nil {
		// Cache miss, query database
		fmt.Println("Retrieving data from the database for id:", id)
		row := db.QueryRow("SELECT name, description, price FROM products WHERE id = $1", id)
		var name, description string
		var price float64
		err := row.Scan(&name, &description, &price)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		product := map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"price":       price,
		}
		// Serialize product to JSON string for caching
		jsonProduct, err := json.Marshal(product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize product"})
			return
		}
		// Set the product in the cache with a TTL of 15 seconds
		rdb.Set(ctx, id, jsonProduct, 15*time.Second)
		c.JSON(http.StatusOK, product)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error accessing cache"})
	} else {
		// Cache hit, deserialize JSON string to map
		fmt.Println("Retrieved data from Redis for id:", id)
		var cachedProduct map[string]interface{}
		err := json.Unmarshal([]byte(val), &cachedProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deserialize cached product"})
			return
		}
		c.JSON(http.StatusOK, cachedProduct)
	}
}

func main() {
	db = initDB()
	rdb = initRedis()

	router := gin.Default()
	router.GET("/products/:id", getProductByID)
	router.Run(":8080")
}
