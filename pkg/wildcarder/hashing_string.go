//go:build no_hashing
// +build no_hashing

package wildcarder

// QuestionHash is the type of the question stored in the cache.
type QuestionHash string

// AnswerHash is the type of an answer stored in the cache.
type AnswerHash DNSAnswer

// HashQuestion hashes a question and returns a QuestionHash.
func HashQuestion(question string) QuestionHash {
	return QuestionHash(question)
}

// HashAnswer hashes a DNSAnswer and returns a AnswerHash.
func HashAnswer(answer DNSAnswer) AnswerHash {
	return AnswerHash(answer)
}
