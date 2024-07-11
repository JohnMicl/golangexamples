package main

import (
	"context"
	"fmt"
	"net"

	"net/http"
)

type ContextInfo struct {
	Params map[string]any
}

type contextKey struct {
	name string
}

var BirdnextCtxKey = contextKey{"bird"}

func NewContext() context.Context {
	return context.WithValue(context.Background(), BirdnextCtxKey, &ContextInfo{
		Params: make(map[string]any),
	})
}

type Handler interface {
	MidServe(c context.Context, w http.ResponseWriter, r *http.Request)
}

type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func (fn HandlerFunc) MidServe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fn(ctx, w, r)
}

type Birdnest struct {
	middleware []Handler
	mux        http.Handler
}

func New() *Birdnest {
	return &Birdnest{
		mux: http.DefaultServeMux,
	}
}

func (b *Birdnest) Use(handler Handler) {
	b.middleware = append(b.middleware, handler)
}

func (b *Birdnest) UseFunc(handler HandlerFunc) {
	b.Use(handler)
}

func (b *Birdnest) SetMux(mux http.Handler) {
	b.mux = mux
}

func (b *Birdnest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext()

	for _, handler := range b.middleware {
		handler.MidServe(ctx, w, r)
	}

	b.mux.ServeHTTP(w, r.WithContext(ctx))
}

func (b *Birdnest) Run(addr string) error {
	return http.ListenAndServe(addr, b)
}

func main() {

	var connHandler = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Println("exc conn get address")
		ctx.Value(BirdnextCtxKey).(*ContextInfo).Params["LocalAddrContextKey"] = r.Context().Value(http.LocalAddrContextKey)
	}

	var authHandler = func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Println("exc conn get token")
		token := r.URL.Query().Get("token")
		if token == "123456" {
			ctx.Value(BirdnextCtxKey).(*ContextInfo).Params["Vailed"] = true
		}
	}

	var b = New()
	b.UseFunc(connHandler)
	b.UseFunc(authHandler)

	// mux := httprouter.New()
	// mux.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 	ctx := r.Context()
	// 	params := ctx.Value(BirdnextCtxKey).(*ContextInfo).Params

	// 	localaddress := params["LocalAddrContextKey"].(*net.TCPAddr)
	// 	valied := params["Vailed"]

	// 	w.Write([]byte(fmt.Sprintf("hello world, localaddr: %+v, valid: %+v\n", *localaddress, valied)))
	// })
	mux := http.NewServeMux()
	wrap := func(method func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			method(w, r)
		}
	}
	mux.Handle("/", wrap(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("exc real request")

		ctx := r.Context()
		params := ctx.Value(BirdnextCtxKey).(*ContextInfo).Params

		localaddress := params["LocalAddrContextKey"].(*net.TCPAddr)
		valied := params["Vailed"]

		w.Write([]byte(fmt.Sprintf("hello world, localaddr: %+v, valid: %+v\n", *localaddress, valied)))
	}))

	b.SetMux(mux)

	b.Run(":8888")
}
