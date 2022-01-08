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

func IsLower(c rune) bool {
	return c >= 'a' && c <= 'z'
}

func IsUpper(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func ToLower(c rune) rune {
	return c - 'A' + 'a'
}

func ToUpper(c rune) rune {
	return c - 'a' + 'A'
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
		matches = strings.TrimSpace(flag.Arg(0))

		if len(matches) != 5 {
			log.Fatal("5 letter word")
		}
	}

	fmt.Println(match(matches))
}

func match(m string) (res []string) {
	misplaced := ""

	for _, c := range m {
		if IsLower(c) {
			misplaced += string(ToUpper(c))
		}
	}

word_loop:
	for _, w := range words {
		if strings.ContainsAny(w, skips) {
			continue
		}

		if !(ContainsAll(w, contains) && ContainsAll(w, misplaced)) {
			continue
		}

		if w != "" {
			compare := ""

			for i, c := range m {
				sc := string(c)

				switch {
				case IsUpper(c):
					compare += string(w[i])

				case IsLower(c):
					if rune(w[i]) == ToUpper(c) {
						continue word_loop
					}
					compare += sc

				default: // symbol
					compare += sc
				}
			}

			if compare != m {
				continue
			}
		}

		res = append(res, w)
	}

	return
}
