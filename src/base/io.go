package base

import (
	"bufio"
	"io"
	"os"
)

func CopySkip(inPath string, outPath string, line int) error {
	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	r := bufio.NewReader(inFile)
	w := bufio.NewWriter(outFile)
	defer w.Flush()
	lineNum := 0
	for {
		buf, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}
		lineNum++
		if lineNum != line {
			_, err := w.Write(buf)
			if err != nil {
				return err
			}
		}
	}
}
