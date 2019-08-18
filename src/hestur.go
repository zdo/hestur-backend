package main

import "./hestur"

func main() {
	game := hestur.NewGame()
	server := hestur.NewServer(game)
	server.Serve()
}
