package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	//go:embed gordle.list
	worddata string

	words    []string
	skips    string
	contains string
)

func ContainsAll(s, chars string) bool {
	for _, c := range chars {
		if !strings.Contains(s, string(c)) {
			return false
		}
	}

	return true
}

func main() {
	wordfile := flag.String("words", "", "external file with words")
	flag.StringVar(&skips, "skip", "", "skip words containing any of these letters")
	flag.StringVar(&contains, "contain", "", "accept only words that contain all these letters")
	flag.Parse()

	skips = strings.ToUpper(strings.TrimSpace(skips))
	contains = strings.ToUpper(strings.TrimSpace(contains))

	var reader io.Reader

	if *wordfile != "" {
		f, err := os.Open(*wordfile)
		if err != nil {
			log.Fatal(err)
		}

		reader = f
	} else {
		reader = strings.NewReader(worddata)
	}

	words = make([]string, 0, 2048)

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		w := strings.ToUpper(strings.TrimSpace(scanner.Text()))

		if len(w) == 5 {
			words = append(words, w)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if closer, ok := reader.(io.Closer); ok {
		closer.Close()
	}

	matches := ""

	if flag.NArg() > 0 {
		matches = strings.ToUpper(strings.TrimSpace(flag.Arg(0)))

		if len(matches) != 5 {
			log.Fatal("5 letter word")
		}
	}

	for _, w := range words {
		if strings.ContainsAny(w, skips) {
			continue
		}

		if !ContainsAll(w, contains) {
			continue
		}

		if matches != "" {
			compare := ""

			for i, c := range matches {
				if c == '-' || c == '#' || c == '*' {
					compare += string(c)
				} else {
					compare += string(w[i])
				}
			}

			if compare != matches {
				continue
			}
		}

		fmt.Println(w)
	}
}
