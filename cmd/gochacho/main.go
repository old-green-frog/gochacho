package main

import "gochacho/internal/server"

func main() {
	s := server.New()
	s.Run()
}
