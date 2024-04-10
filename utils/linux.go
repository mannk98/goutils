package utils

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

func SysGetHomeDir() (home string) {
	home, err := os.UserHomeDir()
	if err == nil {
		return home
	} else {
		return ""
	}
}

func SysGetUsername() string {
	if user, err := user.Current(); err == nil {
		return user.Username
	} else {
		return ""
	}
}

func SysGetUserId() string {
	if user, err := user.Current(); err == nil {
		return user.Uid
	} else {
		return ""
	}
}

func SysGetGroupId() string {
	if user, err := user.Current(); err == nil {
		return user.Gid
	} else {
		return ""
	}
}

func ArgsGet(index int, args []string) string {
	if len(args) > index {
		return args[index]
	} else {
		return ""
	}
}

/*
Linux cmd
*/

func Touch(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	return file.Close()
	//	return nil
}

func Cat(filename ...string) error {
	for _, finame := range filename {
		f, err := os.Open(finame)
		finamecp := f
		if err != nil {
			return err
		}
		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			return err
		}
		defer finamecp.Close()
	}
	return nil
}

func Which(filePath, PATH string) (string, bool) {
	execbasename := filepath.Base(filePath)
	for _, val := range strings.Split(PATH, string(os.PathListSeparator)) {
		if PathIsFile(filepath.Join(val, execbasename)) {
			return filepath.Join(val, execbasename), true
		}
	}
	return "", false
}

func Cp(sourceFile, destinationFile string) bool {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Println(err)
		return false
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		log.Println("Error creating", destinationFile)
		log.Println(err)
		return false
	}
	return true
}

/*
Container
*/
func matchContainerIDWithHostname(lines string) string {
	hostname := os.Getenv("HOSTNAME")
	re := regexp.MustCompilePOSIX("^[[:alnum:]]{12}$")

	if re.MatchString(hostname) {
		regex := fmt.Sprintf("(%s[[:alnum:]]{52})", hostname)

		return matchContainerID(regex, lines)
	}
	return ""
}

func matchContainerID(regex, lines string) string {
	// Attempt to detect if we're on a line from a /proc/<pid>/mountinfo file and modify the regexp accordingly
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt section 3.5
	re := regexp.MustCompilePOSIX("^[0-9]+ [0-9]+ [0-9]+:[0-9]+ /")
	if re.MatchString(lines) {
		regex = fmt.Sprintf("containers/%v", regex)
	}

	re = regexp.MustCompilePOSIX(regex)
	if re.MatchString(lines) {
		submatches := re.FindStringSubmatch(string(lines))
		containerID := submatches[1]

		return containerID
	}
	return ""
}

/*
GetCurrentContainerID attempts to extract the current container ID from the provided file paths. (call inside container)

		If no files paths are provided, it will default to /proc/1/cpuset, /proc/self/cgroup and /proc/self/mountinfo.
	 It attempts to match the HOSTNAME first then use the fallback method, and returns with the first valid match.
*/
func GetCurrentContainerID(filepaths ...string) (id string) {
	if len(filepaths) == 0 {
		filepaths = []string{"/proc/1/cpuset", "/proc/self/cgroup", "/proc/self/mountinfo"}
	}

	// We try to match a 64 character hex string starting with the hostname first
	for _, filepath := range filepaths {
		file, err := os.Open(filepath)
		if err != nil {
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			_, lines, err := bufio.ScanLines([]byte(scanner.Text()), true)
			if err == nil {
				strLines := string(lines)
				if id = matchContainerIDWithHostname(strLines); len(id) == 64 {
					return
				}
			}
		}
	}

	// If we didn't get any ID that matches the hostname, fall back to matching the first 64 character hex string
	for _, filepath := range filepaths {
		file, err := os.Open(filepath)
		if err != nil {
			continue
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			_, lines, err := bufio.ScanLines([]byte(scanner.Text()), true)
			if err == nil {
				strLines := string(lines)
				if id = matchContainerID("([[:alnum:]]{64})", strLines); len(id) == 64 {
					return
				}
			}
		}
	}

	return
}

func IsContainer() bool {
	return len(GetCurrentContainerID()) != 0
	// return gogrep.FileIsMatchLiteralLine("/proc/self/cgroup", "docker") || gogrep.FileIsMatchLiteralLine("/proc/self/cgroup", "lxc")
}
