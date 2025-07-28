package gl

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strings"
)

# FROM Go Playground

type skipTillReader struct {
	rdr   *bufio.Reader
	delim []byte
	found bool
}

func newSkipTillReader(reader io.Reader, delim []byte) *skipTillReader {
	return &skipTillReader{
		rdr:   bufio.NewReader(reader),
		delim: delim,
		found: false,
	}
}

func (str *skipTillReader) Read(p []byte) (n int, err error) {
	if str.found {
		return str.rdr.Read(p)
	} else {
		// search byte by byte for the delimiter
	outer:
		for {
			for i := range str.delim {
				var c byte
				c, err = str.rdr.ReadByte()
				if err != nil {
					n = 0
					return
				}
				// doens't match so start over
				if str.delim[i] != c {
					continue outer
				}
			}
			str.found = true
			// we read the delimiter so add it back
			str.rdr = bufio.NewReader(io.MultiReader(bytes.NewReader(str.delim), str.rdr))
			return str.Read(p)
		}
	}
}

type readTillReader struct {
	rdr   *bufio.Reader
	delim []byte
	found bool
}

func newReadTillReader(reader io.Reader, delim []byte) *readTillReader {
	return &readTillReader{
		rdr:   bufio.NewReader(reader),
		delim: delim,
		found: false,
	}
}

func (rtr *readTillReader) Read(p []byte) (n int, err error) {
	if rtr.found {
		return 0, io.EOF
	} else {
	outer:
		for n < len(p) {
			for i := range rtr.delim {
				var c byte
				c, err = rtr.rdr.ReadByte()
				if err != nil && n > 0 {
					err = nil
					return
				} else if err != nil {
					return
				}
				p[n] = c
				n++
				if rtr.delim[i] != c {
					continue outer
				}
			}
			rtr.found = true
			break
		}
		if n == 0 {
			err = io.EOF
		}
		return
	}
}

// Capture string with specified header and footer from a full string sample.
// Returns the captured string and any error from ioutil.Readall.
func Capture(full string, head string, footer string) (string, error) {
	r := strings.NewReader(full)
	str := newSkipTillReader(r, []byte(head))
	rtr := newReadTillReader(str, []byte(footer))
	bs, err := ioutil.ReadAll(rtr)
	return string(bs), err
}
