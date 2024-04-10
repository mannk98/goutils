package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// CheckErr prints the msg with the prefix 'Error:' and exits with error code 1. If the msg is nil, it does nothing.
func CheckErr(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}

// WriteStringAndCheck writes a string into a buffer, and checks if the error is not nil.
func WriteStringAndCheck(b io.StringWriter, s string) {
	_, err := b.WriteString(s)
	CheckErr(err)
}

// return C:/path/to/file.txt on windows(fileToURL("C:\\path\\to\\file.txt")) and path/to/file.txt on Linux (fileToURL("/path/to/file.txt")).
func FileToURL(in string) string {
	i := strings.Split(in, string(filepath.Separator))
	return path.Join(i...)
}

func SliceStringRemoveDuplicates(slice []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range slice {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}
