package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	args := os.Args

	// Statement to output the edited file
	outputFile, content := ValidateArgs(args)

	res := keywordInstance(string(content))

	err := os.WriteFile(outputFile, []byte(res), 0o664)
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateArgs(args []string) (string, []byte) {
	// Handling if the user didn't inputed the input file or output file
	if len(args) <= 1 {
		fmt.Println("You need to input a filename")
		os.Exit(1)
	}
	if len(args) <= 2 {
		fmt.Println("You forgot the output filename")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Reading content and handling error occurrence
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	return outputFile, content
}

func keywordInstance(text string) string {
	re := regexp.MustCompile(\((cap|up|low)(,\s*(\d+))?\))
	matches := re.FindAllStringSubmatch(text, -1)

	if len(matches) != 0 {
		for _, match := range matches {
			keyword := match[0]
			numAlpha := match[3]
			numInt := 0

			if numAlpha != "" {
				numInt, _ = strconv.Atoi(numAlpha)
			}
			text = ProcessKeyword(text, keyword, numInt)
		}
	}
	return text
}

func ProcessKeyword(text, keyword string, num int) string {
	words := strings.Fields(text)
	wordsSize := len(words)
	newWords := []string{}

	if num > 0 {
		newWords = NumIsGreaterThan0(keyword, num, words)
	} else {
		switch keyword {

		case "(low)":
			for i := 0; i < wordsSize; i++ {
				if words[i] == keyword && i != 0 {
					newWords[i-1] = strings.ToLower(words[i-1])
					i++
				}
				newWords = append(newWords, words[i])
			}
		case "(up)":
			for i := 0; i < wordsSize; i++ {
				if words[i] == keyword && i != 0 {
					newWords[i-1] = strings.ToUpper(words[i-1])
					i++
				}
				newWords = append(newWords, words[i])
			}
		case "(cap)":
			for i := 0; i < wordsSize; i++ {
				if words[i] == keyword && i != 0 {
					lWord := newWords[i-1]
					lWord = strings.ToUpper(string(lWord[0])) + strings.ToLower(lWord[1:])
					newWords[i-1] = lWord
					i++
				}
				newWords = append(newWords, words[i])
			}
		}
	}
	return strings.Join(newWords, " ")
}

func NumIsGreaterThan0(keyword string, num int, words []string) []string {
	newWords := []string(words)

	for i := 0; i < len(words)-1; i++ {
		combined := words[i] + " " + words[i+1]
		if combined == keyword {

			for j := i - num; j <= i; j++ {
				switch {
				case strings.Contains(keyword, "(low,"):
					newWords[j] = strings.ToLower(words[j])
				case strings.Contains(keyword, "(up,"):
					newWords[j] = strings.ToUpper(words[j])
				case strings.Contains(keyword, "(cap,"):
					word := words[j]
					newWords[j] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
				}
			}
			i++
		}
		newWords[i] = words[i]
	}
	return newWords
}