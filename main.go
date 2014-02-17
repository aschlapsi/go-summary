package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func splitContentToSentences(content string) []string {
	s := strings.Replace(content, "\n", ". ", -1)
	return strings.Split(s, ". ")
}

func splitContentToParagraphs(content string) []string {
	return strings.Split(content, "\n\n")
}

func convertStringSlice(orig []string) []interface{} {
	s := make([]interface{}, len(orig))
	for i, v := range orig {
		s[i] = interface{}(v)
	}
	return s
}

func sentencesIntersection(sentence1 string, sentence2 string) float64 {
	s1 := convertStringSlice(strings.Split(sentence1, " "))
	s2 := convertStringSlice(strings.Split(sentence2, " "))
	set1 := mapset.NewSetFromSlice(s1)
	set2 := mapset.NewSetFromSlice(s2)
	intersection := set1.Intersect(set2)

	if intersection.Cardinality() == 0 {
		return 0
	}

	intersectionSize := intersection.Cardinality()
	set1Size := set1.Cardinality()
	set2Size := set2.Cardinality()
	return float64(intersectionSize) / ((float64(set1Size) + float64(set2Size)) / 2)
}

func formatSentence(sentence string) string {
	re, _ := regexp.Compile(`\W+`)
	return string(re.ReplaceAll([]byte(sentence), []byte("")))
}

func getSentencesRanks(content string) map[string]float64 {
	sentences := splitContentToSentences(content)

	n := len(sentences)
	values := make(map[int]map[int]float64)
	for i := 0; i < n; i++ {
		values[i] = make(map[int]float64)
		for j := 0; j < n; j++ {
			values[i][j] = sentencesIntersection(sentences[i], sentences[j])
		}
	}

	sentencesMap := make(map[string]float64)
	for i := 0; i < n; i++ {
		var score float64 = 0
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			score += values[i][j]
		}
		sentencesMap[formatSentence(sentences[i])] = score
	}

	return sentencesMap
}

func getBestSentence(paragraph string, sentencesMap map[string]float64) string {
	sentences := splitContentToSentences(paragraph)

	if len(sentences) < 2 {
		return ""
	}

	bestSentence := ""
	maxValue := 0.0
	for _, sentence := range sentences {
		strippedSentence := formatSentence(sentence)
		if len(strippedSentence) > 0 {
			if sentencesMap[strippedSentence] > maxValue {
				maxValue = sentencesMap[strippedSentence]
				bestSentence = sentence
			}
		}
	}

	return bestSentence
}

func getSummary(content string) string {
	sentencesMap := getSentencesRanks(content)
	paragraphs := splitContentToParagraphs(content)

	summaryBuffer := bytes.NewBufferString("")
	summaryBuffer.WriteString("\n")

	for _, paragraph := range paragraphs {
		sentence := strings.TrimSpace(getBestSentence(paragraph, sentencesMap))
		if len(sentence) > 0 {
			summaryBuffer.WriteString(sentence)
			summaryBuffer.WriteString("\n")
		}
	}

	return summaryBuffer.String()
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalln("Input file missing")
	}

	input := flag.Arg(0)

	fmt.Printf("Processing %s...\n", input)

	content, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("File '%s' could not be opened.", input)
	}

	summary := getSummary(string(content))

	fmt.Println(summary)
	fmt.Println()
	fmt.Printf("Original length %d\n", len(content))
	fmt.Printf("Summary length %d\n", len(summary))
	fmt.Printf("Summary ratio: %.2f%%\n", (100 - (100 * (float64(len(summary)) / (float64(len(content)))))))
}
