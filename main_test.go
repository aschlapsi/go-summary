package main

import (
	"testing"
)

func TestSplitContentToSentences(t *testing.T) {
	content := `
  This is the first sentence. This is the second sentence.
  This is the third sentence.
  `
	expect := []string{
		"",
		"  This is the first sentence",
		"This is the second sentence.",
		"  This is the third sentence.",
		"  ",
	}
	sentences := splitContentToSentences(content)
	if len(sentences) != len(expect) {
		t.Errorf("Expected %d sentences, but got %d", len(expect), len(sentences))
	}
	for i := 0; i < len(expect); i++ {
		if sentences[i] != expect[i] {
			t.Errorf("Got '%s' as result #%d, but expected '%s'.", sentences[i], i, expect[i])
		}
	}
}
