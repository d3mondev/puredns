package wildcarder

import (
	"sync"
)

type answerCache struct {
	cache map[AnswerHash][]string
	mu    sync.Mutex
}

func newAnswerCache() *answerCache {
	cache := answerCache{}
	cache.cache = make(map[AnswerHash][]string)

	return &cache
}

// add adds DNS answers to the list and associate a root domain to them.
func (c *answerCache) add(root string, answers []DNSAnswer) {
	c.mu.Lock()
	answerHashes := []AnswerHash{}
	for _, answer := range answers {
		answerHashes = append(answerHashes, HashAnswer(answer))
	}
	c.mu.Unlock()

	c.addHash(root, answerHashes)
}

// addHash adds DNS answer hashes to the list and associate a root domain to them.
func (c *answerCache) addHash(root string, answers []AnswerHash) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, answer := range answers {
		if _, ok := c.cache[answer]; !ok {
			c.cache[answer] = []string{}
		}

		found := false
		for _, r := range c.cache[answer] {
			if root == r {
				found = true
				break
			}
		}

		if !found {
			c.cache[answer] = append(c.cache[answer], root)
		}
	}
}

// find returns the root domains associated with DNS answers.
func (c *answerCache) find(answers []DNSAnswer) []string {
	c.mu.Lock()

	answerHashes := []AnswerHash{}
	for _, answer := range answers {
		answerHashes = append(answerHashes, HashAnswer(answer))
	}

	c.mu.Unlock()

	return c.findHash(answerHashes)
}

// findHash returns the root domains associated with DNS answer hashes.
func (c *answerCache) findHash(answers []AnswerHash) []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, answer := range answers {
		if roots, ok := c.cache[answer]; ok {
			return roots
		}
	}

	return []string{}
}

// count returns the number of answers in the cache.
func (c *answerCache) count() int {
	var count int

	for _, answers := range c.cache {
		count += len(answers)
	}

	return count
}
