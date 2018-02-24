package main

import (
	"log"
	"strconv"

	"github.com/g-harel/edelweiss/internal/database"
	"github.com/g-harel/edelweiss/internal/database/model"
	"github.com/g-harel/edelweiss/internal/session"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	ss, err := session.NewStore("localhost:6379", "password123")
	if err != nil {
		log.Fatal(err)
	}
	sm := session.NewManager(ss)

	d, err := database.New(`
		host=localhost
		port=5432
		user=postgres
		password=password123
		dbname=edelweiss
		sslmode=disable
	`)
	if err != nil {
		log.Fatal(err)
	}
	_ = model.New(d)

	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	router.GET("/", func(c *gin.Context) {
		s, err := sm.Load(c)
		if err != nil {
			panic(err)
		}

		sessionID, err := s.Get("id")
		if err != nil {
			panic(err)
		}

		visits, _ := s.Get("visits")
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

	router.Run()
}
