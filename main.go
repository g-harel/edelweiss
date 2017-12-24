package main

import (
	"log"
	"os"
	"time"

	"github.com/g-harel/edelweiss/internal/session"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	sm, err := session.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	defer sm.Close()
	r.Use(sm.Middleware)

	r.GET("/", func(c *gin.Context) {
		s := sm.Load(c)

		sessionID := s.Get("id")

		c.JSON(200, gin.H{
			"message": sessionID,
		})
	})

	r.GET("/e", func(c *gin.Context) {
		c.Redirect(301, "/")
		go func() {
			time.Sleep(time.Millisecond * 200)
			go os.Exit(0)
		}()
	})

	r.Run()
}
