package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

var flagInputFile string
var flagOutputFile string
var flagDir string

func init() {
	flag.StringVar(&flagInputFile, "i", "", "input file")
	flag.StringVar(&flagOutputFile, "o", "", "output file")
	flag.StringVar(&flagDir, "d", "", "directory")
}

var kFunc = []rune{'f', 'u', 'n', 'c'}
var kType = []rune{'t', 'y', 'p', 'e'}

func match(src []rune, target []rune) bool {
	if len(src) < len(target) {
		return false
	}
	for i := 0; i < len(target); i++ {
		if src[i] != target[i] {
			return false
		}
	}
	return true
}

/*
func() {           func() {
}                  }
func() {    --->   // add line break
}                  func() {
                   }
//some comment -->  // some comment      add space after //
 */
func format(input, output string) {
	content, _ := ioutil.ReadFile(input)
	rs := bytes.Runes(content)
	builder := strings.Builder{}
	length := len(rs)
	for i := 0; i < length; i++ {
		r := rs[i]
		builder.WriteRune(r)
		if r == rune('/') {
			i++
			if i >= length {
				break
			}
			r = rs[i]
			builder.WriteRune(r)

			if r == rune('/') {
				i++
				if i >= length {
					break
				}

				r = rs[i]
				if !unicode.IsSpace(r) && r != rune('"') {
					builder.WriteRune(rune(' '))
				}
				builder.WriteRune(r)
			}
		}

		if r == rune('}') {
			tmp := make([]rune, 0)
			i++
			if i >= length {
				break
			}

			for ;i < length; i++ {
				if unicode.IsSpace(rs[i]) {
					tmp = append(tmp, rs[i])
				} else {
					break
				}
			}

			if i >= length {
				for _, r := range tmp {
					builder.WriteRune(r)
				}
				break
			}

			word := make([]rune, 0)
			for ; i < length; i++ {
				if rs[i] > 'a' && rs[i] < 'z' {
					word = append(word, rs[i])
					if len(word) == 4 { // "func"
						break
					}
				} else {
					i--
					break
				}
			}

			// fmt.Println("}", string(tmp), string(word))
			if match(word, kFunc) || match(word, kType) {
				builder.WriteRune('\n')
				builder.WriteRune('\n')
			} else {
				for _, r := range tmp {
					builder.WriteRune(r)
				}
			}

			for _, r := range word {
				builder.WriteRune(r)
			}
		}
	}
	tmpFile := input + ".tmp"
	ioutil.WriteFile(tmpFile, []byte(builder.String()), os.ModePerm)
	os.Rename(tmpFile, output)
}

func gofmt(input string) {
	cmd := exec.Command("gofmt", "-l", "-w", input)
	cmd.Start()
	cmd.Wait()
}

func main() {
	flag.Parse()
	if flagInputFile != "" {
		fmt.Println("process: ", flagInputFile)
		outputFile := flagOutputFile
		if outputFile == "" {
			outputFile = flagInputFile
		}
		format(flagInputFile, outputFile)
		gofmt(flagInputFile)
	}

	if flagDir != "" {
		err := filepath.Walk(flagDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || filepath.Ext(path) != ".go" {
				return nil
			}
			fmt.Println("process: ", path)
			format(path, path)
			gofmt(path)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}
