package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/gobs/sortedmap"
)

var (
	//go:embed gordle.list
	worddata string

	words []string
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

func IsLetter(c rune) bool {
	return IsLower(c) || IsUpper(c)
}

func ToLower(c rune) rune {
	return c - 'A' + 'a'
}

func ToUpper(c rune) rune {
	return c - 'a' + 'A'
}

func main() {
	var skips, contains string

	wordfile := flag.String("words", "", "external file with words")
	flag.StringVar(&skips, "skip", "", "skip words containing any of these letters")
	flag.StringVar(&contains, "contain", "", "accept only words that contain all these letters")
	fsort := flag.Bool("sort", true, "sort by frequency")
	all := flag.Bool("all", false, "show all words (answer and allowed) - false: only show valid answer words")
	asList := flag.Bool("list", false, "list mode: the first word is the initial one, the remaining show the results")
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

		if len(w) == 0 && *all == false { // first batch is valid answer words, 2nd batch is allowed words
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if closer, ok := reader.(io.Closer); ok {
		closer.Close()
	}

	sort.Strings(words)

	if *asList {
		processList(flag.Args(), *fsort)
	} else if flag.NArg() > 0 {
		processWord(flag.Arg(0), skips, contains, *fsort)
	} else {
		printList(words, *fsort)
	}
}

func processWord(word, skips, contains string, sort bool) {
	word = strings.TrimSpace(word)
	if len(word) != 5 {
		log.Fatalf("expected 5 letter word: %q", word)
	}

	list := match(word, skips, contains)
	if len(list) == 0 {
		fmt.Println("No matches")
		return
	}

	printList(list, sort)
}

func processList(words []string, sort bool) {
	var word, skips string

	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w) != 5 {
			log.Fatalf("expected 5 letter word: %q", w)
		}

		for p, c := range w {
			if !IsLetter(c) {
				skips += word[p : p+1]
			}
		}

		word = w
	}

	processWord(word, skips, "", sort)
}

func printList(list []string, sort bool) {
	if !sort {
		fmt.Println(list)
		return
	}

	// sort by frequency

	byletter := map[string][]string{}
	bycount := map[string]int{}

	for _, w := range list {
		k := w[0:1]

		byletter[k] = append(byletter[k], w)
	}

	for k, v := range byletter {
		bycount[k] = len(v)
	}

	for _, kv := range sortedmap.AsSortedByValue(bycount, false) {
		fmt.Println(kv.Key, kv.Value, byletter[kv.Key])
	}
}

func match(m, skips, contains string) (res []string) {
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
