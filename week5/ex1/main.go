package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

/*
	1. Stores session for users logging in.
	2. Limits access to /ping to allow only one user at a time (with a 5-second delay).
	3. Counts requests to /ping since the server started.
	4. Implements rate limiting of two requests to /ping per user in 60 seconds.
	5. Returns the top 10 users who called /ping the most.
	6. Uses HyperLogLog to track unique users calling /ping.
*/

/*
When a user logs in, a session is created and stored with a unique session ID, which is sent to the user's browser as a cookie.
Each time the user makes a request, their browser sends this session cookie to the server.
This cookie allows the server to know which session belongs to which user.
*/

var ctx = context.Background()
var redisClient *redis.Client
var pingMutex sync.Mutex
var pingCount int

const (
	rateLimitKey     = "rate_limit"
	topUsersKey      = "top_users"
	hllKey           = "ping_hyperloglog"
	sessionUsername  = "username"
	maxRequests      = 2
	rateLimitSeconds = 60
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func initRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis-pass",
		// DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to redis successfully")
	return client
}

func loginHandler(c *gin.Context) {
	session := sessions.Default(c)
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid input"})
		return
	}

	// Respond to the client
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Username and password cannot be empty"})
		return
	}

	// mock userlogin
	if strings.Contains(req.Username, "test") && req.Password == "123456" {
		session.Set(sessionUsername, req.Username)
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
	}
}

func pingHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(sessionUsername) // return interface{}
	if username == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Rate limiting logic
	userKey := fmt.Sprintf("%s:%v", rateLimitKey, username)
	limitCount, _ := redisClient.Get(ctx, userKey).Int()
	if limitCount >= maxRequests {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Rate limit exceeded"})
		return
	}

	// Use a mutex to ensure that only 1 person calls this at a time
	pingMutex.Lock()
	defer pingMutex.Unlock()

	pingCount++

	// track unique users with HyperLogLog
	redisClient.PFAdd(ctx, hllKey, username)

	// increment request counter for user(for top tracking)
	redisClient.ZIncrBy(ctx, topUsersKey, 1, username.(string))

	// set the rate limit in redis (expire in 60s)
	redisClient.Incr(ctx, userKey)
	redisClient.Expire(ctx, userKey, time.Duration(rateLimitSeconds)*time.Second)

	// sleep 5s
	time.Sleep(5 * time.Second)
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// HyperLogLog: count unique users who called /ping
func countHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping_count": pingCount})
}

func topUsersHandler(c *gin.Context) {
	topUsers, _ := redisClient.ZRevRangeWithScores(ctx, topUsersKey, 0, 9).Result() // top 10
	result := make([]gin.H, len(topUsers))                                          // make a key-value pairs with len of topUsers
	for i, user := range topUsers {
		result[i] = gin.H{"username": user.Member, "count": user.Score}
	}
	c.JSON(http.StatusOK, result)
}

func authMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(sessionUsername)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// Initialize Redis client
	redisClient = initRedisClient()

	// Initialize Redis session store
	store, err := redisStore.NewStore(10, "tcp", "localhost:6379", "redis-pass", []byte("secret"))

	if err != nil {
		log.Fatalf("Failed to create redis store %v", err)
	}

	r.Use(sessions.Sessions("mysession", store))

	r.POST("/login", loginHandler)

	authorized := r.Group("/")
	authorized.Use(authMiddleWare())
	{
		authorized.GET("/ping", pingHandler)
		authorized.GET("/count", countHandler)
		authorized.GET("/top", topUsersHandler)
	}

	r.Run(":8000")
}
