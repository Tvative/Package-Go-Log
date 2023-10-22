// Initialize Package
//
// Copyright (c) 2023 Tvative
// All Rights Reserved
//
// Use of this source code is governed by
// certain licenses found in the LICENSE file

// Package GoLog provides a logging mechanism for managing log instances
//
// Go Log is a Go library for flexible and customizable logging. It provides
// an easy way to log messages to both the terminal and log files. You can
// also add color to your terminal log messages for better readability
package GoLog

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// A LogInstance is a struct that holds information about logging/
type LogInstance struct {
	// A LogDestination is the file where the log will be written/
	logDestination *os.File
}

const (
	ColorDefault string = "\x1b[0;0m"  // A ColorDefault represents the ANSI escape sequence for resetting the text color to the default
	ColorRed     string = "\x1b[31;1m" // A ColorRed represents the ANSI escape sequence for setting text color to red
	ColorYellow  string = "\x1b[33;1m" // A ColorYellow represents the ANSI escape sequence for setting text color to yellow
)

const (
	MessageNormal  string = " [ INFO ] " // A MessageNormal represents a normal message identifier
	MessageFatal   string = " [ ERRO ] " // A MessageFatal represents a fatal error message identifier
	MessageWarning string = " [ WARN ] " // A MessageWarning represents a warning message identifier
)

// An Initialize the log data with the provided file destination
// It opens the file specified by fileDestination and prepares it for writing
// If the file cannot be opened, it returns false along with an error message
func Initialize(logDestination string) *LogInstance {
	fileDescriptor, openError := os.Create(logDestination)

	if openError != nil {
		return nil
	}

	return &LogInstance{fileDescriptor}
}

// A printOutPut Print writes the log message to the specified output destinations
func (logInstance *LogInstance) printOutPut(needFileOutput bool, needTerminalOutput bool,
	needTerminalColoredOutput bool, messageType string,
	jsonContent map[string]interface{}, messageContent ...interface{}) {
	var messagePrefix string

	// Generate message prefix

	getTime := time.Now()
	generateLongTime := getTime.Format("2006-01-02 15:04:05")
	generatedTimeMillSeconds := getTime.Nanosecond() / 1e6
	generatedTimeNanoSeconds := getTime.Nanosecond()

	var generatedTime string = generateLongTime + ":" +
		strconv.Itoa(generatedTimeMillSeconds) + ":" +
		strconv.Itoa(generatedTimeNanoSeconds)

	messagePrefix = generatedTime + messageType

	// Print to the file

	if needFileOutput {
		fmt.Fprint(logInstance.logDestination, messagePrefix)
		fmt.Fprint(logInstance.logDestination, messageContent...)

		if jsonContent != nil {
			logInstance.generateJSON(true, false, jsonContent)
		}

		fmt.Fprintln(logInstance.logDestination)
	}

	// Print to the terminal

	if needTerminalOutput && needTerminalColoredOutput {
		var colorCode string

		switch messageType {
		case MessageNormal:
			colorCode = ColorDefault

		case MessageFatal:
			colorCode = ColorRed

		case MessageWarning:
			colorCode = ColorYellow
		}

		fmt.Print(colorCode, messagePrefix)
		fmt.Print(messageContent...)

		if jsonContent != nil {
			logInstance.generateJSON(false, true, jsonContent)
		}

		fmt.Println(ColorDefault)
	} else if needTerminalOutput {
		fmt.Print(messagePrefix)
		fmt.Print(messageContent...)

		if jsonContent != nil {
			logInstance.generateJSON(false, true, jsonContent)
		}

		fmt.Println()
	}

	// Exit if fatal

	if messageType == MessageFatal {
		os.Exit(1)
	}
}

// A generateJSON Generate and print JSON content
func (logInstance *LogInstance) generateJSON(needFileOutPut bool, needTerminalOutput bool,
	jsonData map[string]interface{}) {

	if needFileOutPut {
		fmt.Fprint(logInstance.logDestination, " [")
	}

	if needTerminalOutput {
		fmt.Print(" [")
	}

	for jsonKey, jsonValue := range jsonData {
		if needFileOutPut {
			fmt.Fprint(logInstance.logDestination, " (", jsonKey, ": ", jsonValue, ")")
		}

		if needTerminalOutput {
			fmt.Print(" (", jsonKey, ": ", jsonValue, ")")
		}
	}

	if needFileOutPut {
		fmt.Fprint(logInstance.logDestination, " ]")
	}

	if needTerminalOutput {
		fmt.Print(" ]")
	}
}
