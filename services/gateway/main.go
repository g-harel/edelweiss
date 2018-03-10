package main

import (
	"log"
	"strconv"

	"github.com/g-harel/edelweiss/services/gateway/database"
	"github.com/g-harel/edelweiss/services/gateway/database/model"
	"github.com/g-harel/edelweiss/services/gateway/session"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	ss, err := session.NewStore("redis:6379", "password123")
	if err != nil {
		log.Fatal(err)
	}
	sm := session.NewManager(ss)

	d, err := database.New(`
		host=psql
		port=5432
		user=postgres
		password=password123
		dbname=edelweiss
		sslmode=disable
	`)
	if err != nil {
		log.Fatal(err)
	}
	m := model.New(d)

	usr, err := m.Users.Add("test@gmail.com", "tomato123")
	if err != nil {
		panic(err)
	}
	_, err = m.Users.Authenticate(usr.Email, "tomato123")
	if err != nil {
		panic(err)
	}
	err = m.Users.ChangeHash(usr.Email, "tomato123", "orange123")
	if err != nil {
		panic(err)
	}
	_, err = m.Users.Authenticate(usr.Email, "orange123")
	if err != nil {
		panic(err)
	}
	err = m.Users.ChangeVerified(usr.UUID, true)
	if err != nil {
		panic(err)
	}

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
