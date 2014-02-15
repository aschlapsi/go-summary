package main

import (
	"bytes"
	"fmt"
	"github.com/deckarep/golang-set"
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

func getSummary(title string, content string, sentencesMap map[string]float64) string {
	paragraphs := splitContentToParagraphs(content)

	summaryBuffer := bytes.NewBufferString(strings.TrimSpace(title))
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
	title := "Swayy is a beautiful new dashboard for discovering and curating online content [Invites]"
	content := `
  Lior Degani, the Co-Founder and head of Marketing of Swayy, pinged me last week when I was in California to tell me about his startup and give me beta access. I heard his pitch and was skeptical. I was also tired, cranky and missing my kids – so my frame of mind wasn’t the most positive.

  I went into Swayy to check it out, and when it asked for access to my Twitter and permission to tweet from my account, all I could think was, “If this thing spams my Twitter account I am going to bitch-slap him all over the Internet.” Fortunately that thought stayed in my head, and not out of my mouth.

  One week later, I’m totally addicted to Swayy and glad I said nothing about the spam (it doesn’t send out spam tweets but I liked the line too much to not use it for this article). I pinged Lior on Facebook with a request for a beta access code for TNW readers. I also asked how soon can I write about it. It’s that good. Seriously. I use every content curation service online. It really is That Good.

  What is Swayy? It’s like Percolate and LinkedIn recommended articles, mixed with trending keywords for the topics you find interesting, combined with an analytics dashboard that shows the trends of what you do and how people react to it. I like it for the simplicity and accuracy of the content curation. Everything I’m actually interested in reading is in one place – I don’t have to skip from another major tech blog over to Harvard Business Review then hop over to another major tech or business blog. It’s all in there. And it has saved me So Much Time



  After I decided that I trusted the service, I added my Facebook and LinkedIn accounts. The content just got That Much Better. I can share from the service itself, but I generally prefer reading the actual post first – so I end up sharing it from the main link, using Swayy more as a service for discovery.

  I’m also finding myself checking out trending keywords more often (more often than never, which is how often I do it on Twitter.com).



  The analytics side isn’t as interesting for me right now, but that could be due to the fact that I’ve barely been online since I came back from the US last weekend. The graphs also haven’t given me any particularly special insights as I can’t see which post got the actual feedback on the graph side (however there are numbers on the Timeline side.) This is a Beta though, and new features are being added and improved daily. I’m sure this is on the list. As they say, if you aren’t launching with something you’re embarrassed by, you’ve waited too long to launch.

  It was the suggested content that impressed me the most. The articles really are spot on – which is why I pinged Lior again to ask a few questions:

  How do you choose the articles listed on the site? Is there an algorithm involved? And is there any IP?

  Yes, we’re in the process of filing a patent for it. But basically the system works with a Natural Language Processing Engine. Actually, there are several parts for the content matching, but besides analyzing what topics the articles are talking about, we have machine learning algorithms that match you to the relevant suggested stuff. For example, if you shared an article about Zuck that got a good reaction from your followers, we might offer you another one about Kevin Systrom (just a simple example).

  Who came up with the idea for Swayy, and why? And what’s your business model?

  Our business model is a subscription model for extra social accounts (extra Facebook / Twitter, etc) and team collaboration.

  The idea was born from our day-to-day need to be active on social media, look for the best content to share with our followers, grow them, and measure what content works best.

  Who is on the team?

  Ohad Frankfurt is the CEO, Shlomi Babluki is the CTO and Oz Katz does Product and Engineering, and I [Lior Degani] do Marketing. The four of us are the founders. Oz and I were in 8200 [an elite Israeli army unit] together. Emily Engelson does Community Management and Graphic Design.

  If you use Percolate or read LinkedIn’s recommended posts I think you’ll love Swayy.

  ➤ Want to try Swayy out without having to wait? Go to this secret URL and enter the promotion code thenextweb . The first 300 people to use the code will get access.

  Image credit: Thinkstock
  `
	sentencesMap := getSentencesRanks(content)
	summary := getSummary(title, content, sentencesMap)

	fmt.Println(summary)
	fmt.Println()
	fmt.Printf("Original length %d\n", len(title)+len(content))
	fmt.Printf("Summary length %d\n", len(summary))
	fmt.Printf("Summary ratio: %.2f%%\n", (100 - (100 * (float64(len(summary)) / (float64(len(title) + len(content)))))))
}
