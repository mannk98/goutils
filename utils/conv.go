package utils

import (
	"bufio"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func Int64ToBytes(number int64) []byte {
	big := new(big.Int)
	big.SetInt64(number)
	return big.Bytes()
}

func Int64ToString(number int64) string {
	big := new(big.Int)
	big.SetInt64(number)
	return big.String()
}

func String2Int64(str string) int64 {
	int64Value, _ := strconv.ParseInt(str, 10, 64)
	return int64Value
}

/* return []string, sperate by endline */
func String2lines(str string) []string {
	scanner := bufio.NewScanner(strings.NewReader(str))

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}
	}

	return lines
}

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Reader2lines(f)
}

func Reader2lines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// get n field of string seprate by delimiter string
func StringGetFirstFields(str string, delimiter string, firstNField uint16) string {
	// Split the string using the delimiter
	fields := strings.Split(str, delimiter)

	// Specify the number of fields you want to take from the beginning
	return strings.Join(fields[:firstNField], delimiter)
}
