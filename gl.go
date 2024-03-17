package gl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type RunArgs struct {
	Exe     string
	Args    []string
	Dir     string
	Env     []string
	Stdin   []byte
	Timeout int
}

type panicT struct {
	msg  string
	code int
}

// Interface to execute the given `RunArgs` through `exec.Command`.
// The first return value is a boolean, true indicates success, false otherwise.
// Second value is the standard output of the command.
// Third value is the standard error of the command.
// Fourth value is error string from Run.
func (a RunArgs) Run() (bool, string, string, string) {
	var r bool = true
	/* #nosec G204 */
	cmd := exec.Command(a.Exe, a.Args...)
	if a.Dir != "" {
		cmd.Dir = a.Dir
	}
	if a.Env != nil || len(a.Env) > 0 {
		cmd.Env = append(os.Environ(), a.Env...)
	}
	if a.Stdin != nil || len(a.Stdin) > 0 {
		cmd.Stdin = bytes.NewBuffer(a.Stdin)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var errorStr string
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	var err error
	if a.Timeout > 0 {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				r = false
				errorStr = sig.String()
			}
		}()
		err = cmd.Start()
		if err != nil {
			r = false
			errorStr = err.Error()
		} else {
			timer := time.AfterFunc(time.Duration(a.Timeout)*time.Second, func() {
				err = cmd.Process.Kill()
				if err != nil {
					r = false
					errorStr = err.Error()
				}
			})
			err = cmd.Wait()
			signal.Stop(c)
			if err != nil {
				r = false
				errorStr = err.Error()
			}
			timer.Stop()
		}
	} else {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				r = false
				errorStr = sig.String()
			}
		}()
		err = cmd.Run()
		signal.Stop(c)
		if err != nil {
			r = false
			errorStr = err.Error()
		}
	}
	return r, stdout.String(), stderr.String(), errorStr
}

func IsFile(p string) bool {
	info, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsDir(p string) bool {
	if fi, err := os.Stat(p); err == nil {
		if fi.IsDir() {
			return true
		}
	}
	return false
}

// Returns a function for simple directory or file check.
// StatPath("directory") for directories.
// StatPath() or StatPath("whatever") for files.
// The function returns boolean `true` on successfully check, `false` otherwise.
func StatPath(f string) func(string) bool {
	switch f {
	case "directory":
		return func(i string) bool {
			if fi, err := os.Stat(i); err == nil {
				if fi.IsDir() {
					return true
				}
			}
			return false
		}
	default:
		return func(i string) bool {
			info, err := os.Stat(i)
			if os.IsNotExist(err) {
				return false
			}
			return !info.IsDir()
		}
	}
}

// Returns a function for walking a path for files.
// Files are read and then contents are written to a strings.Builder pointer.
func PathWalker(sh *strings.Builder) func(string, fs.DirEntry, error) error {
	isFile := StatPath("file")
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		/* #nosec G304 */
		if isFile(path) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			str, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			sh.WriteString(string(str)) // length of string and nil err ignored
			err = file.Close()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Reads a file `path` then returns the contents as a string.
// Always returns a string value.
// An empty string "" is returned for nonexistent or unreadable files.
func FileRead(path string) string {
	isFile := StatPath("file")
	/* #nosec G304 */
	if isFile(path) {
		file, err := os.Open(path)
		if err != nil {
			log.Panic(err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Panic(err)
			}
		}()
		str, err := io.ReadAll(file)
		if err != nil {
			log.Panic(err)
		}
		return string(str)
	} else {
		return ""
	}
}

func FileLines(path string) []string {
	isFile := StatPath("file")
	/* #nosec G304 */
	if isFile(path) {
		file, err := os.Open(path)
		if err != nil {
			log.Panic(err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var text []string
		for scanner.Scan() {
			text = append(text, scanner.Text())
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Panic(err)
			}
		}()
		return text
	} else {
		return nil
	}
}

func FileGlob(path string) ([]string, error) {
	p := ""
	for _, r := range path {
		if unicode.IsLetter(r) {
			p += fmt.Sprintf("[%c%c]", unicode.ToLower(r), unicode.ToUpper(r))
		} else {
			p += string(r)
		}
	}
	return filepath.Glob(p)
}

// Insert string argument #2 into index `i` of first argument `a`.
func InsertStr(a []string, b string, i int) []string {
	a = append(a, "")
	copy(a[i+1:], a[i:]) // number of elements copied ignored
	a[i] = b
	return a
}

// Prefix string `s` with pipes "|".
// Used to "prettify" command line output.
// Returns new string.
func PipeStr(prefix string, char string, str string) string {
	replacement := fmt.Sprintf("\n %s %s ", prefix, char)
	str = strings.Replace(str, "\n", replacement, -1)
	return fmt.Sprintf(" %s %s\n %s %s %s", prefix, char, prefix, char, str)
}

// Writes the string `s` to the file `path`.
// It returns any error encountered, nil otherwise.
func StringToFile(path string, s string) error {
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		cerr := fo.Close()
		if err == nil {
			err = cerr
		}
	}()
	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}
	return nil
}

func RecoverPanic() {
	if rec := recover(); rec != nil {
		err := rec.(panicT)
		fmt.Fprintln(os.Stderr, err.msg)
		os.Exit(err.code)
	}
}

func Assert(e error, s string) {
	if e != nil {
		panic(panicT{msg: PipeStr(s, fmt.Sprintf("%s\n", e.Error()), ">"), code: 255})
	}
}

func Bug(s string) {
	panic(panicT{msg: fmt.Sprintf("BUG: %s", s), code: 255})
}

func Panic(s string) {
	panic(panicT{msg: fmt.Sprintf("FATAL: %s", s), code: 1})
}

func Panicf(f string, a ...interface{}) {
	panic(panicT{msg: fmt.Sprintf("FATAL: "+f, a...), code: 1})
}
