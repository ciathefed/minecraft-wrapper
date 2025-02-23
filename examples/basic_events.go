package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	wrapper "github.com/ciathefed/minecraft-wrapper"
)

// In this example we are mimicking the "seed" commands, where when a player
// says "seed" in-game, we are going to capture that message from the GameEvent
// channel and call the Seed() function from the wrapper and have the wrapper
// broadcast the return seed value to all players with the following message:
// "[Server] The world seed is: 9785468184"

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	wpr, err := wrapper.NewDefaultWrapper("server.jar", "1024M", "1024M")
	if err != nil {
		log.Fatalf("failed to create wrapper: %v", err)
	}
	if err := wpr.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	defer wpr.Stop()

	for {
		select {
		case e := <-wpr.GameEvents():
			if e.String() == "player-say" {
				switch e.Data["player_message"] {
				case "seed":
					seed, err := wpr.Seed()
					if err != nil {
						log.Println("err getting seed: ", err)
						continue
					}
					wpr.Say(fmt.Sprintf("The world seed is: %d", seed))
				default:
					log.Println(e.String(), e.Data)
				}
			}
		case <-c:
			wpr.Kill()
			os.Exit(1)
		}
	}
}
