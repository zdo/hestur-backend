package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"./hestur"
)

func main() {
	cpuprofile := "/tmp/hesturprofile"
	f, err := os.Create(cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)

		pprof.StopCPUProfile()
		os.Exit(0)
	}()

	game := hestur.NewGame()

	for i := 0; i < 30; i++ {
		dummy := hestur.NewDummyCharacter()
		dummy.Position.Y = hestur.Coordinate(i + 2)
		game.RegisterCharacter(dummy)
	}

	server := hestur.NewServer(&game)
	server.Serve()
}
