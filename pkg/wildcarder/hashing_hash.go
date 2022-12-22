//go:build !no_hashing
// +build !no_hashing

package wildcarder

import (
	"hash/maphash"

	"github.com/d3mondev/resolvermt"
)

var hashSeed maphash.Seed = maphash.MakeSeed()

// QuestionHash is the type of the question stored in the cache.
type QuestionHash uint64

// AnswerHash is the type of an answer stored in the cache.
type AnswerHash struct {
	Type resolvermt.RRtype
	Hash uint64
}

// HashQuestion hashes a question and returns a QuestionHash.
func HashQuestion(question string) QuestionHash {
	var hasher maphash.Hash
	hasher.SetSeed(hashSeed)
	hasher.WriteString(question)

	return QuestionHash(hasher.Sum64())
}

// HashAnswer hashes a DNSAnswer and returns a AnswerHash.
func HashAnswer(answer DNSAnswer) AnswerHash {
	var hasher maphash.Hash
	hasher.SetSeed(hashSeed)
	hasher.WriteString(answer.Answer)

	return AnswerHash{
		Type: answer.Type,
		Hash: hasher.Sum64(),
	}
}
