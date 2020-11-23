package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	// formatter := render.New(render.Options{
	// 	IndentJSON: true,
	// })
	formatter := render.New(render.Options{
		Directory:  "templates",
		Extensions: []string{".html"},
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrive working directory")
		} else {
			webRoot = root
			fmt.Println(root)
		}
	}
	mx.HandleFunc("/login", login(formatter))
	mx.HandleFunc("/user", login(formatter))
	mx.HandleFunc("/error", login(formatter))
	mx.HandleFunc("/api/test", apiTestHandler(formatter)).Methods("GET")
	mx.HandleFunc("/index", homeHandler(formatter))
	mx.PathPrefix("/").Handler(http.FileServer(http.Dir(webRoot + "/assets/")))
}
