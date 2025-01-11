package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	args := os.Args

	// Statement to output the edited file
	outputFile, content := validateArgs(args)

	res := processText(string(content))

	err := os.WriteFile(outputFile, []byte(res), 0o664)
	if err != nil {
		log.Fatalf("Failed to create the new file: %v", err)
	}
}

func validateArgs(args []string) (string, []byte) {
	// Handling if the user didn't inputed the input file or output file
	if len(args) < 2 {
		fmt.Println("You forgot the input filename")
		os.Exit(1)
	}
	if len(args) < 3 {
		fmt.Println("You forgot the output filename")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if !(strings.HasSuffix(inputFile, ".txt")) || !(strings.HasSuffix(outputFile, ".txt")) {
		fmt.Println("One of the files has incorrect format (.txt)")
		os.Exit(1)
	}

	// Reading content and handling error occurrence
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	return outputFile, content
}

func processText(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		line = processFlag(line)
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func processFlag(text string) string {
	r := regexp.MustCompile(`(?:^|\s+)\((cap|up|low|hex|bin)(,\s+([+-]?\d+))?(?:\s+|$)?\)`)
	matches := r.FindAllStringSubmatch(text, -1)

	if len(matches) != 0 {
		for _, match := range matches {
			keyword := strings.Trim(match[0], " ")
			numAlpha := match[3]
			numInt := 0

			if numAlpha != "" {
				numInt, _ = strconv.Atoi(numAlpha)
			}
			text = processKeyword(text, keyword, numInt)
		}
	}
	text = deleteExtraSpaces(text)
	text = processPunctuation(text)
	return text
}

func processPunctuation(text string) string {
	// Remove spaces from before and after the punctuation
	text = regexp.MustCompile(`\s*([.,!?;:])\s*`).ReplaceAllString(text, "$1")
	// Add space between the punctuation and the next word
	text = regexp.MustCompile(`([.,!?;:])([^.,!?;:])`).ReplaceAllString(text, "$1 $2")
	text = processQuotes(text)
	text = processAa(text)
	return text
}

func processKeyword(text, keyword string, num int) string {
	words := strings.Fields(text)
	wordsSize := len(words)
	newWords := []string{}

	sIndex := 0
	if wordsSize > 0 && words[0] == keyword {
		sIndex = 1
	}

	for i := sIndex; i < wordsSize; i++ {
		if i < wordsSize-1 && num != 0 && keyword == words[i]+" "+words[i+1] {
			if num < 0 {
				num = -num
			}
			s := i - num
			if s < 0 {
				s = 0
			}
			for j := s; j < i && j < len(newWords); j++ {

				switch {
				case strings.Contains(keyword, "(low,"):
					newWords[j] = toLower(newWords[j])
				case strings.Contains(keyword, "(up,"):
					newWords[j] = toUpper(newWords[j])
				case strings.Contains(keyword, "(cap,"):
					newWords[j] = toCapital(newWords[j])
				}
			}
			i++
		} else if words[i] == keyword && i > 0 {
			prevWordIndex := len(newWords) - 1
			if prevWordIndex >= 0 {
				switch keyword {
				case "(low)":
					newWords[prevWordIndex] = toLower(newWords[prevWordIndex])
				case "(cap)":
					newWords[prevWordIndex] = toCapital(newWords[prevWordIndex])
				case "(up)":
					newWords[prevWordIndex] = toUpper(newWords[prevWordIndex])
				case "(hex)":
					newWords[prevWordIndex] = convertBase(newWords[prevWordIndex], 16)
				case "(bin)":
					newWords[prevWordIndex] = convertBase(newWords[prevWordIndex], 2)
				}
			}
			continue
		} else {
			newWords = append(newWords, words[i])
		}
	}

	return strings.Join(newWords, " ")
}

func convertBase(word string, fromBase int) string {
	digitVal, err := strconv.ParseInt(word, fromBase, 64)
	if err != nil {
		fmt.Printf("Error parsing %s: %v\n", word, err)
		return word
	}
	return strconv.FormatInt(digitVal, 10)
}

func toCapital(s string) string {
	runes := []rune(s)
	for i, v := range runes {
		runes[i] = unicode.ToLower(v)
	}

	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func toUpper(s string) string {
	runes := []rune(s)
	for i, v := range runes {
		runes[i] = unicode.ToUpper(v)
	}
	return string(runes)
}

func toLower(s string) string {
	runes := []rune(s)
	for i, v := range runes {
		runes[i] = unicode.ToLower(v)
	}
	return string(runes)
}

func processQuotes(text string) string {
	r := regexp.MustCompile(`(?:^|\w|\s+)'\s*(.+?)\s+'`)
	return r.ReplaceAllString(text, "'$1'")
}

func processAa(text string) string {
	r := regexp.MustCompile(`\b(a|A)\s+([aeiouhAEIOUH]+)`)
	for r.MatchString(text) {
		text = r.ReplaceAllString(text, "${1}n $2")
	}
	return text
}

func deleteExtraSpaces(text string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
}

