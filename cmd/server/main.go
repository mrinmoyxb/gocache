package main

import (
	"fmt"
	"gocache/command"
	"gocache/pkg/cache"
	"time"
)

func main() {
	c := cache.NewCache(time.Second)
	engine := command.NewEngine(c)

	commands := []string{
		"SET name Mrinmoy",
		"GET name",
		"EXISTS name",
		"GET age",
		"DEL name",
		"GET name",
	}

	for _, input := range commands {
		cmd, err := command.Parse(input)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		result, err := engine.Executed(cmd)
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Printf("> %s\n%s\n", input, result)
	}

}
