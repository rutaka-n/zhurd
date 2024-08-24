package main

import (
	"fmt"
	"os"

	"zhurd/internal/config"
)

func main() {
	fmt.Printf("hello, world!\n")
	cfgFile, err := os.Open("./share/config.json.example")
	if err != nil {
		panic(err)
	}
	cfg, err := config.Load(cfgFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("config: %+v\n", cfg)
}
