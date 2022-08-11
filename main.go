package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commandPrompt string = "Enter a command and data: > "

func main() {
	var line string
	var input []string
	var command string
	var data string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	fmt.Printf(commandPrompt)
	for scanner.Scan() {
		line = scanner.Text()
		input = strings.Split(line, " ")
		command = input[0]
		data = ""
		if len(input) > 1 {
			for _, elem := range input[1:] {
				data = data + " " + elem
			}
		}
		if command == "exit" {
			fmt.Println("[Info] Bye!")
			break
		} else {
			fmt.Println(command, data)
			fmt.Printf(commandPrompt)
		}
	}
}
