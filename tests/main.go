package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Lista de status HTTP com pesos — 200 tem mais chance
var statusPool = []int{
	200, 200, 200, 200, 200, // 5x mais chance de ser 200
	400, 401, 403, 404, 408,
	429, 500, 502, 503, 504,
}

func getRandomStatus() int {
	return statusPool[rand.Intn(len(statusPool))]
}

func handler(w http.ResponseWriter, r *http.Request) {
	status := getRandomStatus()
	w.WriteHeader(status)
	fmt.Fprintf(w, "Respondendo com status %d\n", status)
	fmt.Printf("[RECEBIDO] %s %s -> %d\n", r.Method, r.URL.Path, status)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handler)

	port := 8080
	fmt.Printf("Servidor rodando em http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
