package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const maxNotePadSize int = 5
const numCommands int = 4
const listPattern string = "[Info] %d: %s\n"

// Errors
const errorNotePadFull string = "[Error] Notepad is full"

var quit bool = false
var commandPrompt string = "Enter a command and data: > "
var notePad []string = make([]string, 0, maxNotePadSize)

// Known commands
var commands = [numCommands]string{
	"exit",
	"create",
	"list",
	"clear",
}

// Command confirmation phrases
var commandConfirmation = [4]string{
	"[Info] Bye!\n",
	"[OK] The note was successfully created\n",
	"",
	"[OK] All notes were successfully deleted\n",
}

// Commands corresponding actions
var actions = [4]func(str ...string) error{
	actionOnExit,
	actionOnCreate,
	actionOnList,
	actionOnClear,
}

// action on command "exit", fakeArgs are unnecessary
func actionOnExit(fakeArgs ...string) error {
	_ = fakeArgs
	quit = true
	return nil
}

// action on command "create", data - string to be added to notePad
func actionOnCreate(newNote ...string) error {
	if len(newNote) > 0 {
		if len(notePad) == maxNotePadSize {
			return errors.New(errorNotePadFull)
		}
		notePad = append(notePad, newNote[0])
	}
	return nil
}

// action on command "clear"
func actionOnClear(fakeArgs ...string) error {
	_ = fakeArgs
	notePad = nil
	notePad = make([]string, 0, maxNotePadSize)
	return nil
}

// action on command "list"
func actionOnList(fakeArgs ...string) error {
	_ = fakeArgs
	for i, note := range notePad {
		if len(note) > 0 {
			fmt.Printf(listPattern, i+1, note)
		}
	}
	return nil
}

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
			data = strings.Join(input[1:], " ")
		}

		foundCommand := false
		for i, knownCommand := range commands {
			if command == knownCommand {
				err := actions[i](data)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(commandConfirmation[i])
				}
				foundCommand = true
			}
		}
		if !foundCommand {
			fmt.Println(command, data)
		}
		if quit {
			break
		}
		fmt.Printf(commandPrompt)
	}
}
