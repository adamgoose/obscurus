package main

import (
	"fmt"

	"github.com/adamgoose/obscurus/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
