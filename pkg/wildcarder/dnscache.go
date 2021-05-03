package wildcarder

import (
	"sync"
)

// DNSCache represents a cache of DNS queries and answers.
type DNSCache struct {
	mu    sync.Mutex
	cache map[QuestionHash][]AnswerHash
}

// NewDNSCache creates an empty cache.
func NewDNSCache() *DNSCache {
	cache := DNSCache{}
	cache.cache = make(map[QuestionHash][]AnswerHash)
	return &cache
}

// Add adds an answer to the DNS cache. The answer will be appended to the list of
// existing answers for a question if they already exist.
func (c *DNSCache) Add(question string, answers []DNSAnswer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	questionHash := HashQuestion(question)

	if _, ok := c.cache[questionHash]; !ok {
		c.cache[questionHash] = []AnswerHash{}
	}

	for _, answer := range answers {
		answerHash := HashAnswer(answer)

		found := false
		for _, answer := range c.cache[questionHash] {
			if answer == answerHash {
				found = true
				break
			}
		}

		if !found {
			c.cache[questionHash] = append(c.cache[questionHash], answerHash)
		}
	}
}

// Find returns the answers for a given DNS query from the cache.
// The list of answers returned can be empty if the question is in the cache but
// no results were found, or nil if the question is not in the cache.
func (c *DNSCache) Find(question string) []AnswerHash {
	c.mu.Lock()
	defer c.mu.Unlock()

	questionHash := HashQuestion(question)

	if questionMap, ok := c.cache[questionHash]; ok {
		answers := []AnswerHash{}
		answers = append(answers, questionMap...)

		return answers
	}

	return nil
}
