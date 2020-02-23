package main

import (
	"apiTools/cmd"
	"apiTools/modle"
)

func main() {
	err := cmd.InitServer()

	defer modle.CloseIO()

	if err != nil {
		panic(err)
	}
}
