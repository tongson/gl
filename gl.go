package gl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"
)

type RunPrinter func(string)

type RunArg struct {
	Exe     string
	Args    []string
	Dir     string
	Env     []string
	Stdin   []byte
	Timeout int
	Stdout  RunPrinter
	Stderr  RunPrinter
}

type RunOut struct {
	Stdout string
	Stderr string
	Error  string
}

// Interface to execute the given `RunArg` through `exec.Command`.
// The first return value is a boolean, true indicates success, false otherwise.
// Second value is the standard output of the command.
// Third value is the standard error of the command.
// Fourth value is error string from Run.
func (a RunArg) Run() (bool, RunOut) {
	var r bool
	/* #nosec G204 */
	cmd := exec.Command(a.Exe, a.Args...)
	cmd.SysProcAttr = setPgid()
	if a.Dir != "" {
		cmd.Dir = a.Dir
	}
	if a.Env != nil || len(a.Env) > 0 {
		cmd.Env = append(os.Environ(), a.Env...)
	}
	if a.Stdin != nil || len(a.Stdin) > 0 {
		cmd.Stdin = bytes.NewBuffer(a.Stdin)
	}
	var errorStr string
	var err error
	var stdout string
	var stderr string
	done := make(chan struct{})
	soPipe, soErr := cmd.StdoutPipe()
	if soErr != nil {
		return false, RunOut{Stdout: "", Stderr: "", Error: soErr.Error()}
	}
	sePipe, seErr := cmd.StderrPipe()
	if seErr != nil {
		return false, RunOut{Stdout: "", Stderr: "", Error: soErr.Error()}
	}
	var soStr strings.Builder
	var seStr strings.Builder
	soBuf := bufio.NewReader(soPipe)
	seBuf := bufio.NewReader(sePipe)
	go func() {
		soLast := false
		seLast := false
		for {
			line, err := soBuf.ReadString('\n')
			if err == io.EOF {
				soLast = true
			}
			soStr.WriteString(line)
			if a.Stdout != nil {
				a.Stdout(line)
			}
			if soLast {
				break
			}
			if err != nil {
				continue
			}
		}
		for {
			line, err := seBuf.ReadString('\n')
			if err == io.EOF {
				seLast = true
			}
			seStr.WriteString(line)
			if a.Stderr != nil {
				a.Stderr(line)
			}
			if seLast {
				break
			}
			if err != nil {
				continue
			}
		}
		done <- struct{}{}
	}()
	err = cmd.Start()
	if err != nil {
		r = false
		errorStr = err.Error()
	} else {
		var timer *time.Timer
		if (a.Timeout != 0) && (a.Timeout > 0) {
			timer = time.AfterFunc(time.Duration(a.Timeout)*time.Second, func() {
				pid := cmd.Process.Pid
				switch runtime.GOOS {
				case "windows":
					_ = cmd.Process.Kill()
				default:
					sKill(pid)
				}
			})
		}
		<-done
		err = cmd.Wait()
		stdout = soStr.String()
		stderr = seStr.String()
		if err != nil {
			r = false
			errorStr = err.Error()
		} else {
			r = true
		}
		if (a.Timeout != 0) && (a.Timeout > 0) {
			timer.Stop()
		}
	}
	return r, RunOut{Stdout: stdout, Stderr: stderr, Error: errorStr}
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

// Returns a function for walking a path for files.
// Files are read and then contents are written to a strings.Builder pointer.
func PathWalker(sh *strings.Builder) func(string, fs.DirEntry, error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		/* #nosec G304 */
		if IsFile(path) {
			var file *os.File
			var str []byte
			file, oerr := os.Open(path)
			if err != nil {
				return oerr
			}
			defer file.Close()
			str, rerr := io.ReadAll(file)
			if err != nil {
				return rerr
			}
			sh.WriteString(string(str)) // length of string and nil err ignored
			return file.Close()
		}
		return nil
	}
}

// Reads a file `path` then returns the contents as a string.
// Always returns a string value.
// An empty string "" is returned for errors or nonexistent/unreadable files.
func FileRead(path string) string {
	/* #nosec G304 */
	if IsFile(path) {
		var file *os.File
		var str []byte
		file, oerr := os.Open(path)
		if oerr != nil {
			return ""
		}
		defer file.Close()
		str, rerr := io.ReadAll(file)
		if rerr != nil {
			return ""
		}
		_ = file.Close()
		if len(str) > 0 {
			return string(str)
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func FileLines(path string) []string {
	var text []string
	/* #nosec G304 */
	if IsFile(path) {
		var file *os.File
		file, err := os.Open(path)
		if err != nil {
			return text
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			text = append(text, scanner.Text())
		}
		_ = file.Close()
		return text
	} else {
		return text
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

// Prefix string `s` with pipes "│".
// Used to "prettify" command line output.
// Returns new string.
func PipeStr(prefix string, str string) string {
	replacement := fmt.Sprintf("\n %s │ ", prefix)
	str = strings.Replace(str, "\n", replacement, -1)
	return fmt.Sprintf(" %s │\n %s │ %s\n %s │", prefix, prefix, str, prefix)
}

// Writes the string `s` to the file `path`.
// It returns any error encountered, nil otherwise.
func StringToFile(path string, s string) error {
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fo.Close()
	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}
	return fo.Close()
}
