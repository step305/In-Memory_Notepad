package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const numCommands int = 6
const listPattern string = "[Info] %d: %s\n"
const commandPromptRequest string = "Enter a command and data: > "
const maxSizeRequest string = "Enter the maximum number of notes: > "

// Errors and info messages
const errorNotePadFull string = "[Error] Notepad is full"
const errorUnknownCommand string = "[Error] Unknown command"
const errorEmptyNote string = "[Error] Missing note argument"
const errorInvalidInput string = "[Error] Invalid input while getting max notepad size"
const errorNothingUpdate string = "[Error] There is nothing to update"
const errorNothingDelete string = "[Error] There is nothing to delete"
const errorInvalidPosition string = "[Error] Invalid position: %s"
const errorMissingPositionArgument string = "[Error] Missing position argument"
const errorMissingNoteArgument string = "[Error] Missing note argument"
const errorPositionOutOfBoundaries string = "[Error] Position %d is out of the boundaries [1, %d]"
const infoNotePadEmpty string = "[Info] Notepad is empty"
const infoSuccessDelete string = "[OK] The note at position %d was successfully updated\n"
const infoSuccessUpdate string = "[OK] The note at position %d was successfully updated\n"

var quit bool = false
var maxNotePadSize int = 5
var notePad []string = make([]string, 0, maxNotePadSize)

// Known knownCommands
var knownCommands = [numCommands]string{
	"exit",
	"create",
	"list",
	"clear",
	"update",
	"delete",
}

// Command confirmation phrases
var commandConfirmation = [numCommands]string{
	"[Info] Bye!\n",
	"[OK] The note was successfully created\n",
	"",
	"[OK] All notes were successfully deleted\n",
	infoSuccessUpdate,
	infoSuccessDelete,
}

var deleteConfirmationIndex = Where("delete")
var updateConfirmationIndex = Where("update")

// Commands corresponding actions
var actions = [numCommands]func(str ...string) error{
	actionOnExit,
	actionOnCreate,
	actionOnList,
	actionOnClear,
	actionOnUpdate,
	actionOnDelete,
}

func Where(s string) int {
	for i := range knownCommands {
		if knownCommands[i] == s {
			return i
		}
	}
	return -1
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

func getPosition(s string) (int, error) {
	position, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New(fmt.Sprintf(errorInvalidPosition, s))
	}

	if position < 1 || position > maxNotePadSize {
		return 0, errors.New(fmt.Sprintf(errorPositionOutOfBoundaries, position, maxNotePadSize))
	}

	return position, nil
}

// action on command "update
func actionOnUpdate(data ...string) error {
	if len(data) > 0 {
		if len(data[0]) == 0 {
			return errors.New(errorMissingPositionArgument)
		}
		input := strings.Split(data[0], " ")
		if len(input) == 1 {
			return errors.New(errorMissingNoteArgument)
		}
		position, err := getPosition(input[0])
		if err != nil {
			return err
		}

		if position > len(notePad) {
			return errors.New(errorNothingUpdate)
		}

		newData := strings.Join(input[1:], " ")
		notePad[position-1] = newData
		commandConfirmation[updateConfirmationIndex] = fmt.Sprintf(infoSuccessUpdate, position)
	}
	return nil
}

func deleteNote(position int) {
	notePadNew := make([]string, len(notePad)-1, maxNotePadSize)
	notePadNew = append(notePad[:position], notePad[position+1:]...)
	notePad = notePadNew
}

// action on command "delete"
func actionOnDelete(data ...string) error {
	if len(data) > 0 {
		if len(data[0]) == 0 {
			return errors.New(errorMissingPositionArgument)
		}

		input := strings.Split(data[0], " ")

		position, err := getPosition(input[0])
		if err != nil {
			return err
		}

		if position > len(notePad) {
			return errors.New(errorNothingDelete)
		}

		deleteNote(position - 1)
		commandConfirmation[deleteConfirmationIndex] = fmt.Sprintf(infoSuccessDelete, position)
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
