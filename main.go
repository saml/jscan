package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var supportedTypes map[string]int = map[string]int{
	"mta": 1,
}

func main() {
	var mongoURI string
	var port int
	flag.StringVar(&mongoURI, "mongo", envOr("MONGO_URI", "mongodb://localhost/"), "MongoDB URI")
	flag.IntVar(&port, "port", toInt(os.Getenv("HTTP_PORT"), 8888), "Port to listen to for HTTP requests")
	flag.Parse()

	s := dbConnect(mongoURI)
	defer s.Close()

	db := s.DB("jsonbin")
	coll := db.C("submissions")

	r := gin.New()
	r.GET("/submissions/:ty", func(c *gin.Context) {
		ty := c.Param("ty")
		if _, ok := supportedTypes[ty]; !ok {
			c.String(http.StatusNotFound, "Invalid Type")
			return
		}

		s, err := findSubmittedData(coll, ty, 0, 10)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		json.NewEncoder(c.Writer).Encode(map[string]interface{}{
			"a": s,
		})
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Listening to http://localhost:%d", port)
	r.Run(addr)

}
