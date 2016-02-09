package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
	c := db.C("submissions")

	http.HandleFunc("/submissions/", func(w http.ResponseWriter, r *http.Request) {
		ty := strings.Split(r.URL.Path, "/")[2]
		if _, ok := supportedTypes[ty]; !ok {
			http.Error(w, "Invalid Type", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			s, err := findSubmittedData(c, ty, 0, 10)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"a": s,
			})
		case "POST":
			now := time.Now()
			s := SubmittedData{
				Type:        ty,
				CreatedDate: now,
				UpdateDate:  now,
			}
			json.NewEncoder(w).Encode(&s)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Listening to http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))

}
