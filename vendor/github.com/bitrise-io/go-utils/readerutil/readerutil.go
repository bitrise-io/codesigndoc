package readerutil

import "bufio"

// ReadLongLine - an alternative to bufio.Scanner.Scan,
// which can't handle long lines. This function is slower than
// bufio.Scanner.Scan, but can handle arbitrary long lines.
func ReadLongLine(r *bufio.Reader) (string, error) {
	// Do NOT create a `bufio.Reader` inside thise function,
	// get it as an input! (just in case you'd thing about doing a "revision" on this)
	// Creating the `bufio.Reader` here would reset/alter the reader,
	// if it would be created for every line! Not a good idea!

	isPrefix := true
	var err error
	var line, ln []byte

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}
