package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	GOOS       = runtime.GOARCH
	AppName    string
	AppVersion string
	NEWLINE    = "\n"
	Ipv4Regex  = `([0-9]+\.){3}[0-9]+`
)

func init() {
	GOOS := runtime.GOOS
	if GOOS != "windows" {
		NEWLINE = "\r\n"
	} else if GOOS != "darwin" {
		NEWLINE = "\r"
	}
}

/*
Path
*/
func PathBaseName(filePath string) string {
	return filepath.Base(filePath)
}

func PathIsExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func PathIsFile(path string) bool {
	if finfo, err := os.Stat(path); err == nil {
		if !finfo.IsDir() {
			return true
		}
	}
	return false
}

func PathIsDir(path string) bool {
	if finfo, err := os.Stat(path); err == nil {
		if finfo.IsDir() {
			return true
		}
	}
	return false
}

/* add data to PATH , vd: /home/mannk:path */
func PATHJointList(PATH, data string) string {
	//	data = data + string(os.PathSeparator)
	if len(PATH) == 0 {
		return data
	}
	return PATH + string(os.PathListSeparator) + data
	//	filepath.ListSeparator
}

/* remove addpath from path*/
func PATHRemove(PATH, addpath string) string {
	if len(PATH) == 0 {
		return ""
	}
	newpath := ""
	for i, val := range strings.Split(PATH, string(os.PathListSeparator)) {
		if !(val == addpath) {
			if i == 0 {
				newpath = val
			} else {
				newpath = newpath + string(os.PathListSeparator) + val
			}
		}
	}
	return newpath
	//	filepath.ListSeparator
}

func PATHGetEnvPathValue() string {
	for _, pathname := range []string{"PATH", "path"} {
		path := os.Getenv(pathname)
		if len(path) != 0 {
			return path
		}
	}
	return ""
}

/* return PATH as array (:)*/
func PATHArr() []string {
	envs := PATHGetEnvPathValue()
	if len(envs) != 0 {
		return strings.Split(envs, string(os.PathListSeparator))
	}
	return []string{}
}

/* retur path or PATH*/
func PathGetEnvPathKey() string {
	for _, pathname := range []string{"PATH", "path"} {
		path := os.Getenv(pathname)
		if len(path) != 0 {
			return pathname
		}
	}
	return ""
}

/*
File
*/
func FileCreate(fullPath string) error {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		file, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		return file.Close()
	} else {
		return err
	}
}

func FileCreateOverwrite(fullPath string) error {
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	return file.Close()
}

func FileCreateWithContent(fullPath string, data []byte) (bytewrite int, err error) {
	file, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	if err != nil {
		return 0, err
	}
	n, err := file.Write(data)
	return n, err
}

// open and attend to file
func FileOpen2Write(fullPath string) (*os.File, error) {
	// create dir to path if it is not exist
	err := DirCreate(filepath.Dir(fullPath), 0775)
	if err != nil {
		return nil, err
	}
	// create file if it not existed
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = FileCreate(fullPath)
		if err != nil {
			return nil, err
		}
	}

	logf, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	return logf, err
}

func FileReadAll(fullPath string) (string, error) {
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(content), err
}

func FileIsWriteable(path string) (isWritable bool) {
	isWritable = false

	if file, err := os.OpenFile(path, os.O_WRONLY, 0666); err == nil {
		defer file.Close()
		isWritable = true
	} else {
		if os.IsPermission(err) {
			return false
		}
	}
	return
}

/* cp date-time modify from dst file to src file */
func FileCloneDate(dst, src string) bool {
	var err error
	var srcinfo os.FileInfo
	if srcinfo, err = os.Stat(src); err == nil {
		if err = os.Chtimes(dst, srcinfo.ModTime(), srcinfo.ModTime()); err == nil {
			return true
		}
	}
	//	fmt.Errorf("Cannot clone date file ", err)
	return false
}

/* get Md5 string of file */
func FileHashMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

/*
	waitForFile waits for the specified file to exist before returning. If the an

error, other than the file not existing, occurs, the error is returned. If,
after 100 attempts, the file does not exist, an error is returned.
*/
func FileWaitForFileExist(path string, timeoutms int) error {
	if timeoutms < 50 && timeoutms != 0 {
		timeoutms = 50
	}

	for i := 0; i < timeoutms/50; i++ {
		_, err := os.Stat(path)
		if err == nil || !os.IsNotExist(err) {
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("file does not exist: %s", path)
}

/*
	readFile waits for the specified file to contain contents, and then returns

those contents as a string. If an error occurs while reading the file, the
error is returned. If the file has no content after 100 attempts, an error is
returned.
*/
func FileWaitContentsAndRead(path string, timeoutms int) (string, error) {
	if timeoutms < 50 && timeoutms != 0 {
		timeoutms = 50
	}
	for i := 0; i < timeoutms/50; i++ {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		if len(bytes) > 0 {
			return strings.TrimSpace(string(bytes)), err
		}
		time.Sleep(50 * time.Millisecond)
	}
	return "", fmt.Errorf("file is empty: %s", path)
}

/* removeFile removes the specified file. Errors are ignored.*/
func FileRemoveFile(path string) error {
	return os.Remove(path)
}

/* only read filesize, not dir*/
func FileGetSize(filepath string) (int64, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	// get the size
	return fi.Size(), nil
}

/* write content to file if it different with exist file content*/
func FileWriteStringIfChange(pathfile string, contents []byte) (bool, error) {

	oldContents := []byte{}
	if _, err := os.Stat(pathfile); err == nil {
		oldContents, _ = ioutil.ReadFile(pathfile)
	}

	//if bytes.Compare(oldContents, contents) != 0 {
	if !bytes.Equal(oldContents, contents) {
		return true, ioutil.WriteFile(pathfile, contents, 0644)
	} else {
		return false, nil
	}
}

/* cp source to dir. If dest exist, overwrite it  */
func FileCopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file ", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err == nil {
		os.Chmod(dst, sourceFileStat.Mode())
		//os.Chtimes(dst, sourceFileStat.ModTime(), sourceFileStat.ModTime())
		//os.Chown(dst, int(sourceFileStat.Sys().(*syscall.Stat_t).Uid), int(sourceFileStat.Sys().(*syscall.Stat_t).Gid))
		//fmt.Println(int(sourceFileStat.Sys().(*syscall.Stat_t).Uid), int(sourceFileStat.Sys().(*syscall.Stat_t).Gid))
	}
	return nBytes, err
}

/* cp source to dir, keep modify time, owner. If dest exist, overwrite it  */
func FileCopyKeepOwner(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file ", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err == nil {
		os.Chmod(dst, sourceFileStat.Mode())
		os.Chtimes(dst, sourceFileStat.ModTime(), sourceFileStat.ModTime())
		os.Chown(dst, int(sourceFileStat.Sys().(*syscall.Stat_t).Uid), int(sourceFileStat.Sys().(*syscall.Stat_t).Gid))
		//fmt.Println(int(sourceFileStat.Sys().(*syscall.Stat_t).Uid), int(sourceFileStat.Sys().(*syscall.Stat_t).Gid))
	}
	return nBytes, err
}

func FileMove(src, dst string) (int64, error) {
	nbytes, err := FileCopyKeepOwner(src, dst)
	if err != nil {
		return nbytes, err
	}
	FileRemoveFile(src)
	return nbytes, err
}

func FileMoveSameFileSys(src, dst string) error {
	// Get information about the source file
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if the source file is a regular file
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	// Rename the source file to the destination
	if err := os.Rename(src, dst); err != nil {
		return err
	}

	// Set the permissions and modification time of the destination file
	if err := os.Chmod(dst, sourceFileStat.Mode()); err != nil {
		return err
	}
	if err := os.Chtimes(dst, sourceFileStat.ModTime(), sourceFileStat.ModTime()); err != nil {
		return err
	}

	return nil
}

/* insert string at index lines of file, if this line have text, push it one line down. */
func FileInsertStringAtLine(filePath, str string, index int) error {
	NEWLINE := "\n"
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	str = str + NEWLINE //add newline
	scanner := bufio.NewScanner(f)
	lines := ""
	linenum := 0
	inserted := false
	for scanner.Scan() {
		linenum = linenum + 1
		if linenum == index {
			inserted = true
			lines = lines + str
		}
		lines = lines + scanner.Text() + NEWLINE
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if !inserted {
		if index == -1 {
			index = linenum + 1
		}
		for i := linenum + 1; i < index; i++ {
			lines = lines + NEWLINE
		}
		lines = lines + str
	}
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, []byte(lines), info.Mode().Perm())
}

/* create tempdir and return tempfile in tempdir (not create yet)*/
func FileTempCreateInNewTemDir(filename string) string {

	rootdir, err := ioutil.TempDir("", "system")
	if err != nil {
		return ""
	} else {
		//			defer os.RemoveAll(dir)
	}
	return filepath.Join(rootdir, filename)
}

func FileTempCreateInNewTemDirWithContent(filename string, data []byte) string {
	rootdir, err := ioutil.TempDir("", "system")
	if err != nil {
		return ""
	}
	fPath := filepath.Join(rootdir, filename)
	err = os.WriteFile(fPath, data, 0755)
	if err != nil {
		os.RemoveAll(rootdir)
		return ""
	}
	return fPath
}

func walk(filename string, linkDirname string, walkFn filepath.WalkFunc) error {
	symWalkFunc := func(path string, info os.FileInfo, err error) error {

		if fname, err := filepath.Rel(filename, path); err == nil {
			path = filepath.Join(linkDirname, fname)
		} else {
			return err
		}

		if err == nil && info.Mode()&os.ModeSymlink == os.ModeSymlink {
			finalPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			info, err := os.Lstat(finalPath)
			if err != nil {
				return walkFn(path, info, err)
			}
			if info.IsDir() {
				return walk(finalPath, path, walkFn)
			}
		}

		return walkFn(path, info, err)
	}
	return filepath.Walk(filename, symWalkFunc)
}

// Walk extends filepath.Walk to also follow symlinks
func SymWalk(path string, walkFn filepath.WalkFunc) error {
	return walk(path, path, walkFn)
}

/*
Dir
*/
func DirCreate(dirPath string, permission fs.FileMode) error {
	dirFullPath := dirPath
	// if _, err := os.Stat(dirFullPath); err == nil {
	// 	if err := os.RemoveAll(dirFullPath); err != nil {
	// 		//fmt.Println("Error removing existing directory:", err)
	// 		return err
	// 	}
	// }
	/*err := os.MkdirAll(dirFullPath, os.ModePerm)*/
	err := os.MkdirAll(dirFullPath, permission)
	if err != nil {
		//fmt.Println("Error creating the directory:", err)
		return err
	}
	return err
	//fmt.Println("Directory created successfully at:", fullPath)
}

func DirRemove(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	return err
}

/* list all file and dir int level1 of dir*/
func DirAllChild(directory string) (files []string, err error) {
	dirRead, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	dirFiles, err := dirRead.Readdir(0)
	if err != nil {
		return nil, err
	}
	for index := range dirFiles {
		fileHere := dirFiles[index]

		// Get name of file and its full path.
		nameHere := fileHere.Name()
		/*		fmt.Println(nameHere)*/
		/*fullPath := directory + nameHere*/
		fullPath := filepath.Join(directory, nameHere)
		files = append(files, fullPath)
		// Remove the file.
	}
	return files, err
}

/* Remove all content of dir */
func DirRemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Copy dir, change ownner (as user execute), change datetime modify
Notes: Care full with recursive can't overload the File decriptor limit
*/
func DirCopy(srcDir, dstDir string) error {
	_, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("%s %s", srcDir, err)
	}
	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// Read the contents of the source directory
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	// Iterate over the entries in the source directory
	for _, entry := range entries {
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(dstDir, entry.Name())

		// If the entry is a directory, recursively copy it
		if entry.IsDir() {
			if err := DirCopy(src, dst); err != nil {
				return err
			}
		} else {
			// If the entry is a file, copy it
			if _, err := FileCopy(src, dst); err != nil {
				return err
			}
		}
	}
	return nil
}

// Notes: Care full with recursive can't overload the File decriptor limit
func DirCopyKeepOwner(srcDir, dstDir string) error {
	srcStat, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("%s %s", srcDir, err)
	}
	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}
	os.Chown(dstDir, int(srcStat.Sys().(*syscall.Stat_t).Uid), int(srcStat.Sys().(*syscall.Stat_t).Gid))

	// Read the contents of the source directory
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	// Iterate over the entries in the source directory
	for _, entry := range entries {
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(dstDir, entry.Name())

		// If the entry is a directory, recursively copy it
		if entry.IsDir() {
			if err := DirCopyKeepOwner(src, dst); err != nil {
				return err
			}
		} else {
			// If the entry is a file, copy it
			if _, err := FileCopyKeepOwner(src, dst); err != nil {
				return err
			}
		}
	}
	return nil
}

func DirMove(src, dst string) error {
	// Rename the source directory to the destination
	if !PathIsDir(src) {
		return fmt.Errorf("%s is not a directory", src)
	}
	/* 	if err := os.Rename(src, dst); err != nil {
		return err
	} */

	if err := DirCopy(src, dst); err != nil {
		return err
	}
	DirRemove(src)

	return nil
}

func DirMoveKeepOwner(src, dst string) error {
	// Rename the source directory to the destination
	if !PathIsDir(src) {
		return fmt.Errorf("%s is not a directory", src)
	}
	/* 	if err := os.Rename(src, dst); err != nil {
		return err
	} */

	if err := DirCopyKeepOwner(src, dst); err != nil {
		return err
	}
	DirRemove(src)

	return nil
}

func DirIsEmptyorNotexist(directory string) bool {
	dirRead, err := os.Open(directory)
	if err != nil {
		return true
	}
	dirFiles, err := dirRead.Readdir(0)
	if err != nil {
		return true
	}
	if len(dirFiles) > 0 {
		return false
	} else {
		return true
	}
}

// use for both file and dir
func CreateSymLink(src, dest string) error {
	// Create the symbolic link
	err := os.Symlink(src, dest)
	if err != nil {
		return err
	}
	return err
}

// parser shell path include ${shelvariale}
func ParserShellPath(shellPath string) string {
	//shellPath := "${HOSTHOME}/var_log_cloudstack/${HOSTHOME}/${HOME}|/var/log/cloudstack|d|y"
	// Loop through the string by treating it as a slice of bytes
	var shellVars []string
	var tempShellVar string
	var isVariable bool
	for i := 0; i < len(shellPath); i++ {
		// Check if the current character is the start of a variable
		if shellPath[i] == '$' && i+1 < len(shellPath) && shellPath[i+1] == '{' {
			isVariable = true
			tempShellVar = ""
			i++ // Skip the '{' character
			continue
		}

		// Check if the current character is the end of a variable
		if isVariable && shellPath[i] == '}' {
			shellVars = append(shellVars, tempShellVar)
			isVariable = false
			continue
		}

		// Add characters to the temporary variable if inside a variable
		if isVariable {
			tempShellVar += string(shellPath[i])
		}
	}

	shellVars = SliceStringRemoveDuplicates(shellVars)
	//fmt.Println(shellVars)

	for _, v := range shellVars {
		shellVal := os.Getenv(v)
		shellPath = strings.ReplaceAll(shellPath, "$"+"{"+v+"}", shellVal)
	}
	//fmt.Println(shellPath)
	return shellPath
}
