package main

import (
	"errors"
	"strings"
	"strconv"
)

// Check if the string is full of digits.
func isDigit(data string) bool {
	for _, val := range data {
		num, err := strconv.Atoi(string(val))
		if err != nil {
			return false
		}
		if num < 0 && num > 9 {
			return false
		}
	}
	return true
}

// Parses if the date is valid.
func parseDate(data string) error {
	if !isDigit(data[:4]) || !isDigit(data[5:7]) || !isDigit(data[8:]) {
		return errors.New("Not valid incoming data format for Date [" + data + "]")
	}
	if data[4] != '-' || data[7] != '-' {
		return errors.New("Not valid incoming data format for Date [" + data + "]")
	}
	return nil
}

// There are three different formats of incoming string: 2020-04-09, 04:06, 2.
// Date, time, number of comments respectively.
// The function checks if the string is valid for all types.
func checkInput(data string) (int, error) {
	data = strings.TrimSpace(data)
	if len(data) == 10 && (data[4] == '-' || data[7] == '-') {
		if err := parseDate(data); err != nil {
			return -1, err
		} else {
			fullDate = data
			return 0, nil
		}
	} else if len(data) == 5 && data[2] == ':' {
		return 1, nil
	} else if isDigit(data) {
		return 2, nil
	} else {
		return -1, errors.New("Unknown incoming data format [" + data + "]")
	}
}
