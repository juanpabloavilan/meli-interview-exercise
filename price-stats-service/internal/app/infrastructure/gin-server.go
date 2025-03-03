package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	// Good Practice: Set cors
	// Good Practice: Set trace middleware

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"statusCode": http.StatusNotFound, "message": "Page not found"})
	})

	return r
}

func RunHTTPServer(r *gin.Engine, port string) {
	port = fmt.Sprintf(":%s", port)
	err := r.Run(port)
	if err != nil {
		panic(err)
	}
}
