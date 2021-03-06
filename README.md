# `routing`

... lets you decribe the structure of your (restful) API as a struct. The
goal is to then generate an [openapi/swagger](https://swagger.io/) documentation
on the fly.

## Project status

Do not use this thing yet!

## Usage

Setup your server:

```golang
package main

import (
    // ...

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type Server struct {
	listener string
	handler  http.Handler
}

func NewServer(listener string) Server {
	s := Server{listener: listener}
	r := mux.NewRouter().StrictSlash(true)
	routes := s.sitesRoutes()
	routes.Populate(r, "sites")
	s.handler = alice.New().Then(r)
	return s
}

func (s Server) Run() {
	fmt.Printf("Serving at http://%s\nPress CTRL-c to stop...\n", s.listener)
	log.Fatal(http.ListenAndServe(s.listener, s.handler))
}

func (s Server) respond(res http.ResponseWriter, req *http.Request, code int, data interface{}) {
    // ...
}
```

Define your routing:

```golang
package main

import (
	r "github.com/unprofession-al/routing"
)

func (s Server) sitesRoutes() r.Route {
	return r.Route{
		H: r.Handlers{"GET": r.Handler{F: s.SitesHandler, Q: []*r.QueryParam{formatParam}}},
		R: r.Routes{
			"{site}": {
				R: r.Routes{
					"status":  {H: r.Handlers{"GET": r.Handler{F: s.StatusHandler, Q: []*r.QueryParam{formatParam}}}},
					"publish": {H: r.Handlers{"PUT": r.Handler{F: s.PublishHandler, Q: []*r.QueryParam{formatParam}}}},
					"update":  {H: r.Handlers{"PUT": r.Handler{F: s.UpdateHandler, Q: []*r.QueryParam{formatParam}}}},
					"files": {
						H: r.Handlers{"GET": r.Handler{F: s.TreeHandler, Q: []*r.QueryParam{formatParam}}},
						R: r.Routes{
							"*": {
								H: r.Handlers{
									"GET":  {F: s.FileHandler, Q: []*r.QueryParam{mdParam}},
									"POST": {F: s.FileWriteHandler, Q: []*r.QueryParam{mdParam}},
								},
							},
						},
					},
				},
			},
		},
	}
}

var formatParam = &r.QueryParam{
	N:    "f",
	D:    "json",
	Desc: "format of the output, can be 'yaml' or 'json'",
}

var mdParam = &r.QueryParam{
	N: "o",
	D: "all",
	Desc: `define if only one part of the markdown file is requested,
can be 'fm' for frontmatter, 'md' for markdown, all for everything`,
}
```

Prepare the handlers:

```golang
package main

import (
	"net/http"
)

func (s Server) SitesHandler(res http.ResponseWriter, req *http.Request) {
	s.respond(res, req, http.StatusOK, "Hello world!")
}

// ...
```

And fire off:

```golang
package main

func main() {
	s := NewServer("127.0.0.1:8080")
	s.Run()
}
```
