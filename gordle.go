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

func IsLetters(w string) bool {
	for _, c := range w {
		if !IsLetter(c) {
			return false
		}
	}

	return true
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

	res := newResults(skips, contains)

	if *asList {
		res.processList(flag.Args())
	} else if flag.NArg() > 0 {
		res.processWord(flag.Arg(0))
	}

	res.printMatches(*fsort)
}

type results struct {
	skips     string
	contains  string // this is only used via command line
	misplaced []string
	correct   []string
}

func newResults(skips, contains string) *results {
	return &results{
		skips:     skips,
		contains:  contains,
		misplaced: make([]string, 5),
		correct:   make([]string, 5),
	}
}

func (r *results) parse(in, out string) {
	for i, c := range out {
		sc := string(c)

		switch {
		case IsLower(c):
			r.misplaced[i] += strings.ToUpper(sc)

		case IsUpper(c):
			r.correct[i] = sc

		default:
			r.skips += string(in[i])
		}
	}
}

func (r *results) matches() (res []string) {
	var contains = r.contains

	for _, c := range r.misplaced {
		if c != "" {
			contains += c
		}
	}

	for _, c := range r.correct {
		if c != "" {
			contains += c
		}
	}

word_loop:
	for _, w := range words {
		if strings.ContainsAny(w, r.skips) {
			continue
		}

		if !ContainsAll(w, contains) {
			continue
		}

		for i, c := range w {
			sc := string(c)

			if strings.Contains(r.misplaced[i], sc) { // correct in the wrong place
				continue word_loop
			}

			if r.correct[i] != "" && sc != r.correct[i] { // correct in the right place
				continue word_loop
			}
		}

		res = append(res, w)
	}

	return
}

func (r *results) printMatches(sort bool) {
	list := r.matches()
	if len(list) == 0 {
		fmt.Println("No matches")
		return
	}

	printList(list, sort)
}

func (r *results) processWord(word string) {
	word = strings.TrimSpace(word)
	if len(word) != 5 {
		log.Fatalf("expected 5 letter word: %q", word)
	}

	r.parse("     ", word)
}

func (r *results) processList(words []string) {
	var word string

	for i, w := range words {
		w = strings.TrimSpace(w)
		if len(w) != 5 {
			log.Fatalf("expected 5 letter word: %q", w)
		}

		if i%2 == 0 {
			if !IsLetters(w) {
				log.Fatalf("expected all letters")
			}

			word = w
			continue
		}

		r.parse(word, w)
	}
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
