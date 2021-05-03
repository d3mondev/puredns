package wildcarder

import (
	"testing"

	"github.com/d3mondev/resolvermt"
	"github.com/stretchr/testify/assert"
)

func TestDNSCacheAdd(t *testing.T) {
	answerA := []DNSAnswer{
		{Type: resolvermt.TypeA, Answer: "127.0.0.1"},
	}

	answerB := []DNSAnswer{
		{Type: resolvermt.TypeA, Answer: "127.0.0.1"},
		{Type: resolvermt.TypeAAAA, Answer: "::1"},
		{Type: resolvermt.TypeCNAME, Answer: "test"},
		{Type: resolvermt.TypeA, Answer: "127.0.0.1"},
	}

	wantA := []AnswerHash{
		HashAnswer(answerA[0]),
	}

	cache := NewDNSCache()

	cache.Add("question", answerA)
	got := cache.Find("question")
	assert.ElementsMatch(t, wantA, got, "element added to internal cache")

	wantB := []AnswerHash{
		HashAnswer(answerA[0]),
		HashAnswer(answerB[1]),
		HashAnswer(answerB[2]),
	}

	cache.Add("question", answerB)
	got = cache.Find("question")
	assert.ElementsMatch(t, wantB, got, "no element duplicated")
}

func TestDNSCacheAddDifferentQuestion(t *testing.T) {
	answerA := []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}}
	answerB := []DNSAnswer{{Type: resolvermt.TypeAAAA, Answer: "::1"}}

	wantA := []AnswerHash{HashAnswer(answerA[0])}
	wantB := []AnswerHash{HashAnswer(answerB[0])}

	cache := NewDNSCache()
	cache.Add("question 1", answerA)
	cache.Add("question 2", answerB)

	got1 := cache.Find("question 1")
	got2 := cache.Find("question 2")

	assert.ElementsMatch(t, wantA, got1)
	assert.ElementsMatch(t, wantB, got2)
}

func TestDNSCacheFind(t *testing.T) {
	answers := []DNSAnswer{
		{Type: resolvermt.TypeA, Answer: "127.0.0.1"},
		{Type: resolvermt.TypeAAAA, Answer: "::1"},
		{Type: resolvermt.TypeCNAME, Answer: "test"},
	}

	hashes := []AnswerHash{
		HashAnswer(answers[0]),
		HashAnswer(answers[1]),
		HashAnswer(answers[2]),
	}

	tests := []struct {
		name         string
		haveAnswers  []DNSAnswer
		haveQuestion string
		want         []AnswerHash
	}{
		{name: "existing question", haveQuestion: "question", haveAnswers: answers, want: hashes},
		{name: "existing question without answers", haveQuestion: "question", haveAnswers: []DNSAnswer{}, want: []AnswerHash{}},
		{name: "inexistent question", haveQuestion: "invalid", haveAnswers: answers, want: nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := NewDNSCache()
			cache.Add("question", test.haveAnswers)

			got := cache.Find(test.haveQuestion)
			assert.ElementsMatch(t, test.want, got)

			if got == nil || test.want == nil {
				assert.Equal(t, test.want, got)
			}
		})
	}
}
