package server

import (
	"net/http"

	"github.com/go-zoo/bone"
)

// http://192.168.1.116:84/TwP6d7o83sSHmjfhKQRzTX/URL Test?url=https://www.baidu.com
func NewRoute() *bone.Mux {
	r := bone.New()
	c := NewController()
	r.Get("/ping", http.HandlerFunc(c.Ping))
	r.Post("/ping", http.HandlerFunc(c.Ping))

	r.Get("/register", http.HandlerFunc(c.Register))
	r.Post("/register", http.HandlerFunc(c.Register))

	r.Get("/:key/:body", http.HandlerFunc(c.Index))
	r.Post("/:key/:body", http.HandlerFunc(c.Index))

	r.Get("/:key/:title/:body", http.HandlerFunc(c.Index))
	r.Post("/:key/:title/:body", http.HandlerFunc(c.Index))

	r.Get("/:key/:category/:title/:body", http.HandlerFunc(c.Index))
	r.Post("/:key/:category/:title/:body", http.HandlerFunc(c.Index))

	r.Get("/stop", http.HandlerFunc(c.Stop))
	r.Post("/stop", http.HandlerFunc(c.Stop))

	return r
}
