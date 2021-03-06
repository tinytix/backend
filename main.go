package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/EveN-FT/backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	mode := os.Getenv("MODE")
	allowOrigins := []string{
		"http://localhost:3000",
	}

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin", "Authorization", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if mode == "prod" {
				for _, allowed := range allowOrigins {
					if allowed == origin {
						return true
					}
				}

				_, err := url.Parse(origin)
				if err != nil {
					return false
				}
			}
			return true
		},
	}))

	fmt.Println("backend")
	fmt.Println("Initializing controllers")

	v1 := r.Group("/api/v1")
	{
		ticketGroup := v1.Group("/ticket")
		{
			ticket := new(controllers.TicketController)
			ticketGroup.POST("/redeem", ticket.Redeem)
			ticketGroup.POST("/transfer", ticket.Transfer)
			ticketGroup.POST("/create", ticket.CreateRedeem)
		}

		eventGroup := v1.Group("/event")
		{
			event := new(controllers.EventController)
			eventGroup.POST("/create", event.Create)
			eventGroup.POST("/list", event.ListEvents)
			eventGroup.POST("/list-by-owner", event.ListEventsByOwner)
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})

	// Add profiling
	if os.Getenv("PROFILING") != "" {
		pprof.Register(r, "debug/pprof")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r.Run(":" + port)
}
