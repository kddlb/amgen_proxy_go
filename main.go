package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS Approximation (Adjust as per your requirements)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Next()
	})

	// Logger Middleware (Basic)
	router.Use(func(c *gin.Context) {
		start := time.Now()
		// Proceed with the request
		c.Next()
		// Log after the request
		fmt.Printf("%s %s - %s\n", c.Request.Method, c.Request.URL.Path, time.Since(start))
	})

	// Routes
	router.GET("/search", searchHandler)
	router.GET("/songs/:id", songsHandler)
	router.GET("/get/:path", getHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8585"
	}
	router.Run(":" + port)
}

func searchHandler(c *gin.Context) {
	response, err := http.Get("https://api.genius.com/search?" + c.Request.URL.RawQuery)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	jsonData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/json", jsonData)
}

func songsHandler(c *gin.Context) {
	response, err := http.Get("https://api.genius.com/songs/" + c.Param("id") + "?" + c.Request.URL.RawQuery)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	jsonData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/json", jsonData)
}

func getHandler(c *gin.Context) {
	url := "https://genius.com/" + c.Param("path")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	req.Header.Set("Cookie", os.Getenv("GENIUS_COOKIE"))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 Edg/122.0.0.0")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, string(body))
}
