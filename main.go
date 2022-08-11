package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const numCommands int = 4
const listPattern string = "[Info] %d: %s\n"
const commandPromptRequest string = "Enter a command and data: > "
const maxSizeRequest string = "Enter the maximum number of notes: > "

// Errors and info messages
const errorNotePadFull string = "[Error] Notepad is full"
const errorUnknownCommand string = "[Error] Unknown command"
const errorEmptyNote string = "[Error] Missing note argument"
const errorInvalidInput string = "[Error] Invalid input while getting max notepad size"
const infoNotePadEmpty string = "[Info] Notepad is empty"

var quit bool = false
var maxNotePadSize int = 5
var notePad []string = make([]string, 0, maxNotePadSize)

// Known knownCommands
var knownCommands = [numCommands]string{
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
		if len(strings.TrimSpace(newNote[0])) == 0 {
			return errors.New(errorEmptyNote)
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
	if len(notePad) == 0 {
		fmt.Println(infoNotePadEmpty)
		return nil
	}
	for i, note := range notePad {
		if len(note) > 0 {
			fmt.Printf(listPattern, i+1, note)
		}
	}
	return nil
}

func getMaxNotePadSize() (int, error) {
	var newSize int
	fmt.Print(maxSizeRequest)
	_, err := fmt.Scanln(&newSize)
	if err != nil {
		return 0, err
	}
	return newSize, nil
}

func getUserInput(scanner *bufio.Scanner) (command string, data string) {
	line := scanner.Text()
	input := strings.Split(line, " ")
	command = input[0]
	data = ""
	if len(input) > 1 {
		data = strings.Join(input[1:], " ")
	}
	return command, data
}

func createScanner() *bufio.Scanner {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	return scanner
}

func printCommandPrompt() {
	fmt.Print(commandPromptRequest)
}

func findCommand(command string) (int, bool) {
	for i, knownCommand := range knownCommands {
		if knownCommand == command {
			return i, true
		}
	}
	return -1, false
}

func executeCommand(i int, data string) {
	err := actions[i](data)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(commandConfirmation[i])
	}
}

func main() {
	var command string
	var data string
	var scanner *bufio.Scanner
	var err error

	scanner = createScanner()
	maxNotePadSize, err = getMaxNotePadSize()
	if err != nil || maxNotePadSize <= 0 {
		fmt.Println(errorInvalidInput)
		return
	}

	printCommandPrompt()
	for scanner.Scan() {
		command, data = getUserInput(scanner)
		commandIndex, found := findCommand(command)
		if found {
			executeCommand(commandIndex, data)
		} else {
			fmt.Println(errorUnknownCommand)
		}
		if quit {
			break
		}
		printCommandPrompt()
	}
}
