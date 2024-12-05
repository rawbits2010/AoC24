// Adds alternative input providing functionality to Advent of Code puzzle solutions.
// Supports input directly from the commandline, a separate file, and from a webpage.
//
// Suggested usage: lines := ReadInput()
package inputhandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ReadInput is a one call method to parse the commandline and return the input data as separate lines.
// On any caught error, it will exit the app with an error text.
func ReadInput() []string {

	inputMethod, paramValue, err := ParseCommandLine()
	if err != nil {
		fmt.Printf("Error while parsing command line: %v\n\n", err)
		fmt.Println("Usage: cmd -[p/f/w] [data/uri]")
		fmt.Println("p - data is provided as a ';' separated value")
		fmt.Println("f - data is in the file pointed to by the provided path")
		fmt.Println("w - data is given by a website pointed to by the provided url")
		os.Exit(int(ErrorCodeParameters))
	}

	var lines []string
	switch inputMethod {
	case InputParameters:
		lines = strings.Split(paramValue, ";")

	case InputFile:
		inputData, err := GetDataFromFile(paramValue)
		if err != nil {
			fmt.Printf("Error while reading from file '%s': %v", paramValue, err)
			os.Exit(int(ErrorCodeFiles))
		}
		lines = strings.Split(strings.TrimSuffix(inputData, "\n"), "\n")

	case InputWebpage:
		inputData, err := GetDataFromWebpage(paramValue)
		if err != nil {
			fmt.Printf("Error while reading from URL '%s': %v", paramValue, err)
			os.Exit(int(ErrorCodeNetwork))
		}
		lines = strings.Split(strings.TrimSuffix(inputData, "\n"), "\n")

	}
	if len(lines) == 0 {
		fmt.Println("Error: no data was given")
		os.Exit(4)
	}

	return lines
}

// ErrorCodes is the suggested application exit codes.
// Your code can use ErrorCodeProcessing just for consistency.
type ErrorCodes int

const (
	ErrorCodeParameters ErrorCodes = 1
	ErrorCodeFiles      ErrorCodes = 2
	ErrorCodeNetwork    ErrorCodes = 3
	ErrorCodeData       ErrorCodes = 4
	ErrorCodeProcessing ErrorCodes = 5
)

// InputMethod is the determined input method from the commandline arguments.
type InputMethod string

const (
	InputInvalid    InputMethod = "InputInvalid"
	InputParameters InputMethod = "InputParameters"
	InputFile       InputMethod = "InputFile"
	InputWebpage    InputMethod = "InputWebpage"
)

// ErrorInvalidParameters returnde by ParseCommandLine when it faild to parse parameters.
var ErrorInvalidParameters = fmt.Errorf("invalid parameters")

// ParseCommandLine is the commandline parser.
// It returns the determined input method, the associated parameter value, or the error if any.
func ParseCommandLine() (InputMethod, string, error) {

	args := os.Args[1:]

	if len(args) < 2 {
		return InputInvalid, "", ErrorInvalidParameters
	}

	switch args[0] {
	case "-p":
		return InputParameters, args[1], nil
	case "-f":
		return InputFile, args[1], nil
	case "-w":
		return InputWebpage, args[1], nil
	}

	return InputInvalid, "", ErrorInvalidParameters
}

// GetDataFromFile will try to open the file at the given path and returns it's contents or an error if any.
func GetDataFromFile(path string) (string, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// GetDataFromWebpage will try to make a GET request to the given URL and returns the downloaded data as text, or an error if any.
// Optionally it reads the contents of the session.txt file if exists and adds it as a "session" cookie to the request.
// (Advent of Code site needs this to identify the current user.)
func GetDataFromWebpage(url string) (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	sessionValue, err := os.ReadFile("session.txt")
	if err == nil {
		req.AddCookie(&http.Cookie{Name: "session", Value: string(sessionValue)})
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
