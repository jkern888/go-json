package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

type JsonArray []interface{}
type JsonObj map[string]interface{}

func main() {
	data, err := ioutil.ReadFile("input.json")

	if err != nil {
		log.Fatal(err)
	}

	input := string(data)
	index := eatWhitespace(input, 0)
	obj, err := readValue(input, &index)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(obj)

}

func readValue(input string, curIndex *int) (interface{}, error) {
	switch input[*curIndex] {
	case '"':
		return readString(input, curIndex)

	case '.', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return readNumber(input, curIndex)

	case '[':
		return readArray(input, curIndex)

	case '{':
		return readMap(input, curIndex)

	case 't', 'f':
		return readBool(input, curIndex)

	case 'n':
		return readNull(input, curIndex)

	default:
		return nil, errors.New(fmt.Sprintf("cannot parse value starting at index %d", *curIndex))
	}

	return nil, nil
}

func readMap(input string, curIndex *int) (JsonObj, error) {
	obj := make(JsonObj)
	startIndex := *curIndex
	*curIndex = eatWhitespace(input, *curIndex+1)

	for input[*curIndex] != '}' {
		if *curIndex >= len(input) {
			return nil, errors.New(fmt.Sprintf("reached end of input without closing object started at index %d", startIndex))
		}

		key, err := readString(input, curIndex)

		if err != nil {
			return nil, err
		}

		*curIndex = eatWhitespace(input, *curIndex+1)
		if input[*curIndex] != ':' {
			return nil, errors.New(fmt.Sprintf("expected ':' at index %d", *curIndex))
		}
		*curIndex = eatWhitespace(input, *curIndex+1)

		value, err := readValue(input, curIndex)

		if err != nil {
			return nil, err
		}

		obj[key] = value

		*curIndex = eatWhitespace(input, *curIndex+1)

		if input[*curIndex] == ',' {
			*curIndex = eatWhitespace(input, *curIndex+1)
		} else if input[*curIndex] != '}' {
			return nil, errors.New(fmt.Sprintf("incorrectly trying to define new key at index %d", *curIndex))
		}
	}

	return obj, nil
}

func readArray(input string, curIndex *int) (JsonArray, error) {
	array := JsonArray{}
	startIndex := *curIndex
	*curIndex = eatWhitespace(input, *curIndex+1)

	for input[*curIndex] != ']' {
		if *curIndex >= len(input) {
			return nil, errors.New(fmt.Sprintf("reached end of input without closing array started at index %d", startIndex))
		}

		value, err := readValue(input, curIndex)

		if err != nil {
			return nil, err
		}

		array = append(array, value)

		*curIndex = eatWhitespace(input, *curIndex+1)

		if input[*curIndex] == ',' {
			*curIndex = eatWhitespace(input, *curIndex+1)
		} else if input[*curIndex] != ']' {
			return nil, errors.New(fmt.Sprintf("incorrectly trying to define new value at index %d", *curIndex))
		}
	}

	return array, nil
}

func readString(input string, curIndex *int) (string, error) {
	if input[*curIndex] != '"' {
		return "", errors.New(fmt.Sprintf("string must start with '\"', at index %d", *curIndex))
	}

	*curIndex++
	keyStart := *curIndex

	for input[*curIndex] != '"' {
		if *curIndex >= len(input) {
			return "", errors.New(fmt.Sprintf("reached end of input without closing string started at index %d", keyStart))
		}

		*curIndex++
	}

	return input[keyStart:*curIndex], nil
}

func readNumber(input string, curIndex *int) (interface{}, error) {
	keyStart := *curIndex

	delimChars := []uint8{',', '\t', '\n', '}', ']', ' '}
	float := false
	run := true

	for run {
		*curIndex++

		for _, char := range delimChars {
			if char == input[*curIndex] {
				run = false
				break
			}
		}

		if *curIndex >= len(input) {
			return "", errors.New(fmt.Sprintf("reached end of input while parsing number, starting at index %d", keyStart))
		}

		if input[*curIndex] == '.' {
			float = true
		}

	}

	// rewind to the last char of the number to leave the index in the right place for the caller
	*curIndex--

	if float {
		return strconv.ParseFloat(input[keyStart:*curIndex+1], 64)
	} else {
		return strconv.ParseInt(input[keyStart:*curIndex+1], 10, 64)
	}
}

func readNull(input string, curIndex *int) (interface{}, error) {
	if *curIndex+4 > len(input) {
		return nil, errors.New("EOF while parsing null")
	} else if input[*curIndex:*curIndex+4] == "null" {
		*curIndex += 3
		return nil, nil
	} else {
		return nil, errors.New(fmt.Sprintf("unexpected null value '%s' at index %d", input[*curIndex:*curIndex+4], *curIndex))
	}
}

func readBool(input string, curIndex *int) (interface{}, error) {
	length := 0
	keyStart := *curIndex

	switch input[*curIndex] {
	case 't':
		length = 4
	case 'f':
		length = 5
	}

	if *curIndex+length > len(input) {
		return nil, errors.New("EOF while parsing bool")
	} else {
		*curIndex += length - 1
		return strconv.ParseBool(input[keyStart:*curIndex])
	}
}

func eatWhitespace(input string, curIndex int) int {
	for curIndex < len(input) {
		char := input[curIndex]

		if char == ' ' || char == '\t' || char == '\n' {
			curIndex++
		} else {
			break
		}
	}

	return curIndex
}
