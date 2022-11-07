package modules

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

var (
	RemoveLineMutex sync.Mutex
	AppendLineMutex sync.Mutex
)

func AppendLine(fileName, line string) {
	AppendLineMutex.Lock()
	defer AppendLineMutex.Unlock()
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	if _, err = file.WriteString(line + "\n"); err != nil {
		panic(err)
	}
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Contains(array []string, s interface{}) bool {
	for _, v := range array {
		if v == s {
			return true
		}
	}

	return false
}

func RemoveFromArray(array []string, s string) []string {
	for i, v := range array {
		if v == s {
			array = append(array[:i], array[i+1:]...)
			return array // Remove this line to remove ALL items in a list thats "s"
		}
	}
	return array
}

func RemoveFromFile(path string, s string) {
	RemoveLineMutex.Lock()
	defer RemoveLineMutex.Unlock()
	lines, err := ReadLines(path)

	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	for _, v := range lines {
		if !strings.Contains(v, s) {
			file.WriteString(v + "\n")
		}
	}
}
