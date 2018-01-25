package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/g-harel/edelweiss/internal/session"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	sm, err := session.NewManager()
	if err != nil {
		log.Fatal(err)
	}
	defer sm.Close()

	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
		sm.Middleware,
	)

	router.GET("/", func(c *gin.Context) {
		s := sm.Load(c)

		sessionID := s.Get("id")

		visits := s.Get("visits")
		v, err := strconv.Atoi(visits)
		if err != nil {
			v = 0
		}
		v++
		s.Set("visits", strconv.Itoa(v))

		c.JSON(200, gin.H{
			"message": sessionID,
			"visits":  v,
		})
	})

	router.GET("/e", func(c *gin.Context) {
		c.Redirect(301, "/")
		go func() {
			time.Sleep(time.Millisecond * 200)
			go os.Exit(0)
		}()
	})

	router.Run()
}
