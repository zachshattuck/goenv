package goenv

import (
	"errors"
	"os"
	"strconv"
)

// Map of commonly-searched characters in readUntil to friendly string names
var charMap = map[byte]string{
	'\n': "newline",
	'\r': "carriage return",
	' ':  "space",
	'\t': "tab",
}

/*
Reads from start of slice until specified character.
If the character is not found, it will return an error instead of the value so far.
*/
func readUntil(buf []byte, char byte, includeChar bool) ([]byte, error) {
	for i, b := range buf {
		if b == char {
			if includeChar {
				return buf[:i+1], nil
			}
			return buf[:i], nil
		}
	}

	// we never found the character they requested
	friendlyName := charMap[char]
	printableChar := strconv.Quote(string(char))
	if friendlyName != "" {
		return nil, errors.New("no " + friendlyName + " (" + printableChar + ")" + " found")
	}

	return nil, errors.New("no instance of " + printableChar + " found")
}

/*
Reads from start of slice until newline character (either \n or \r\n).
If the newline character is not found, it will return an error instead of the value so far.
*/
func readUntilNewline(buf []byte, includeNewline bool) ([]byte, error) {
	for i, b := range buf {
		if b == '\n' {
			if includeNewline {
				return buf[:i+1], nil
			}
			return buf[:i], nil
		}
		// if cr, it may be a Windows-style newline
		if b == '\r' {
			if i+1 < len(buf) && buf[i+1] == '\n' {
				if includeNewline {
					return buf[:i+2], nil
				}
				return buf[:i], nil
			}
		}
	}

	return nil, errors.New("no \"\\n\" or \"\\r\\n\" found")
}

// I should just make a personal package for deserialization of [NAME][SEPARATOR][VALUE][LINETERMINATOR] sequences
func deserAndSetEnvironment(data []byte) error {
	lineNo := 1

	for len(data) > 0 {
		line, err := readUntilNewline(data, true)
		eof := false
		if err != nil {
			line = data
			eof = true
		}

		// if \n or \r\n, continue
		if len(line) == 1 || (len(line) == 2 && line[0] == '\r') {
			if eof {
				break
			}
			data = data[len(line):]
			lineNo++
			continue
		}

		paramName, err := readUntil(line, '=', false)
		if err != nil {
			return errors.New("failed to read param name on line " + strconv.Itoa(lineNo) + ": " + err.Error())
		}

		valueStartIdx := len(paramName) + 1
		var paramValue []byte

		if eof {
			paramValue = line[valueStartIdx:]
		} else {
			paramValue, err = readUntilNewline(line[valueStartIdx:], false)
			if err != nil {
				return errors.New("failed to read param value of  " + string(paramName) + ": " + err.Error())
			}
		}

		os.Setenv(string(paramName), string(paramValue))

		if eof {
			break
		}
		data = data[len(line):]
		lineNo++
	}

	return nil
}

func ProcessEnv() error {
	dat, err := os.ReadFile(".env")
	if err != nil {
		return errors.New("failed to open .env file")
	}

	err = deserAndSetEnvironment(dat)
	if err != nil {
		return err
	}
	return nil
}
