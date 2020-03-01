package main

import (
	"apiTools/cmd"
	"apiTools/modles"
)

func main() {
	err := cmd.InitServer()

	defer modles.CloseIO()

	if err != nil {
		panic(err)
	}
}
