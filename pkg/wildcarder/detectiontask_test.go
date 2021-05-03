package wildcarder

import (
	"sort"
	"testing"

	"github.com/d3mondev/resolvermt"
	"github.com/stretchr/testify/assert"
)

func newStubDetectionTask() (*detectionTask, *fakeResolver) {
	resolver := newFakeResolver()

	ctx := detectionTaskContext{}
	ctx.results = &result{}
	ctx.resolver = resolver
	ctx.dnsCache = NewDNSCache()
	ctx.preCache = NewDNSCache()
	ctx.wildcardCache = newAnswerCache()
	ctx.randomSubs = []string{"random", "random", "random"}
	ctx.queryCount = 3

	task := &detectionTask{
		domain: "",
		ctx:    ctx,
	}

	return task, resolver
}

func TestCheckPrecache(t *testing.T) {
	t.Run("not a wildcard", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		task.ctx.preCache.Add("foo", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})

		got := task.checkPrecache("foo")
		assert.False(t, got)
	})

	t.Run("wildcard found", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		task.ctx.preCache.Add("foo", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
		task.ctx.wildcardCache.add("foo", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})

		got := task.checkPrecache("foo")
		assert.True(t, got)
	})

	t.Run("root doesn't match domain", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		task.ctx.preCache.Add("bar", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
		task.ctx.wildcardCache.add("foo", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})

		got := task.checkPrecache("bar")
		assert.False(t, got)
	})
}

func TestCheckResolve(t *testing.T) {
	t.Run("not a wildcard", func(t *testing.T) {
		task, resolver := newStubDetectionTask()
		resolver.addAnswer("foo", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})

		got := task.checkResolve("foo")
		assert.False(t, got)
	})

	t.Run("wildcard found", func(t *testing.T) {
		task, resolver := newStubDetectionTask()
		task.ctx.wildcardCache.add("foo", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
		resolver.addAnswer("foo", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})

		got := task.checkResolve("foo")
		assert.True(t, got)
	})

	t.Run("root doesn't match domain", func(t *testing.T) {
		task, resolver := newStubDetectionTask()
		task.ctx.wildcardCache.add("bar", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
		resolver.addAnswer("foo", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})

		got := task.checkResolve("foo")
		assert.False(t, got)
	})
}

func TestTestWildcard(t *testing.T) {
	t.Run("topmost domain", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		got := task.testWildcard("test.com")
		assert.Equal(t, "", got)
	})

	t.Run("not a wildcard domain", func(t *testing.T) {
		task, resolver := newStubDetectionTask()
		resolver.addAnswer("www.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
		resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})

		got := task.testWildcard("www.test.com")
		assert.Equal(t, "", got)
	})

	t.Run("wildcard domain", func(t *testing.T) {
		task, resolver := newStubDetectionTask()
		resolver.addAnswer("www.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
		resolver.addAnswer("random.test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})
		resolver.addAnswer("test.com", []DNSAnswer{{Type: resolvermt.TypeA, Answer: "127.0.0.1"}})

		got := task.testWildcard("www.test.com")
		assert.Equal(t, "test.com", got)
	})
}

func TestFindWildcardRoot(t *testing.T) {
	t.Run("topmost domain", func(t *testing.T) {
		task, _ := newStubDetectionTask()

		gotRoot, gotAnswers := task.findWildcardRoot("test.com", []AnswerHash{})

		assert.Equal(t, "test.com", gotRoot)
		assert.Empty(t, gotAnswers)
	})

	t.Run("recurse to topmost", func(t *testing.T) {
		answerA := DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}
		answerB := DNSAnswer{Type: resolvermt.TypeA, Answer: "10.0.0.1"}
		answerC := DNSAnswer{Type: resolvermt.TypeA, Answer: "192.168.0.1"}

		task, resolver := newStubDetectionTask()
		resolver.addAnswer("www.foo.bar.com", []DNSAnswer{answerA})
		resolver.addAnswer("random.foo.bar.com", []DNSAnswer{answerA})
		resolver.addAnswer("foo.bar.com", []DNSAnswer{answerB})
		resolver.addAnswer("random.bar.com", []DNSAnswer{answerB})
		resolver.addAnswer("bar.com", []DNSAnswer{answerC})

		gotRoot, gotAnswers := task.findWildcardRoot("www.foo.bar.com", []AnswerHash{HashAnswer(answerA)})

		assert.Equal(t, "bar.com", gotRoot)
		assert.Equal(t, []AnswerHash{HashAnswer(answerA), HashAnswer(answerB)}, gotAnswers)
	})

	t.Run("answers do not match", func(t *testing.T) {
		answerA := DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}
		answerB := DNSAnswer{Type: resolvermt.TypeA, Answer: "10.0.0.1"}

		task, resolver := newStubDetectionTask()
		resolver.addAnswer("www.foo.bar.com", []DNSAnswer{answerA})
		resolver.addAnswer("random.foo.bar.com", []DNSAnswer{answerA})
		resolver.addAnswer("foo.bar.com", []DNSAnswer{answerB})

		gotRoot, gotAnswers := task.findWildcardRoot("www.foo.bar.com", []AnswerHash{HashAnswer(answerA)})

		assert.Equal(t, "foo.bar.com", gotRoot)
		assert.Equal(t, []AnswerHash{HashAnswer(answerA)}, gotAnswers)
	})
}

func TestResolveRandom(t *testing.T) {
	t.Run("no parent", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		got := task.resolveRandomSubdomains("test.com")
		assert.Nil(t, got)
	})

	t.Run("no record", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		got := task.resolveRandomSubdomains("www.test.com")
		assert.Nil(t, got)
	})

	t.Run("from dns cache", func(t *testing.T) {
		answer := DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}

		task, _ := newStubDetectionTask()
		task.ctx.dnsCache.Add("random.test.com", []DNSAnswer{answer})

		got := task.resolveRandomSubdomains("www.test.com")
		assert.Equal(t, []AnswerHash{HashAnswer(answer)}, got)
	})

	t.Run("from resolver", func(t *testing.T) {
		answer := DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}

		task, resolver := newStubDetectionTask()
		resolver.addAnswer("random.test.com", []DNSAnswer{answer})

		got := task.resolveRandomSubdomains("www.test.com")
		assert.Equal(t, []AnswerHash{HashAnswer(answer)}, got)
	})
}

func TestMakeTestDomains(t *testing.T) {
	t.Run("subsequent calls should be deterministic", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		task.ctx.randomSubs = newRandomSubdomains(5)

		gotA := task.makeTestSubdomains("www.test.com")
		gotB := task.makeTestSubdomains("www.test.com")

		assert.ElementsMatch(t, gotA, gotB)
	})

	t.Run("no element should be equal", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		task.ctx.randomSubs = newRandomSubdomains(5)

		sorted := task.makeTestSubdomains("www.test.com")
		sort.Strings(sorted)

		for i := 0; i < len(sorted)-1; i++ {
			assert.NotEqual(t, sorted[i], sorted[i+1])
		}
	})

	t.Run("no parent should return nil", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		task.ctx.randomSubs = newRandomSubdomains(5)

		got := task.makeTestSubdomains("test.com")
		assert.Nil(t, got)
	})
}

func TestGetParent(t *testing.T) {
	tests := []struct {
		name       string
		haveDomain string
		wantParent string
	}{
		{name: "no parent", haveDomain: "test.com", wantParent: ""},
		{name: "empty domain", haveDomain: "", wantParent: ""},
		{name: "subdomain", haveDomain: "www.test.com", wantParent: "test.com"},
		{name: "long domain", haveDomain: "one.two.three.four.five.test.com", wantParent: "two.three.four.five.test.com"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotParent := getParent(test.haveDomain)

			assert.Equal(t, test.wantParent, gotParent)
		})
	}
}

func TestAnswerMatch(t *testing.T) {
	testAnswerA := HashAnswer(DNSAnswer{
		Type:   resolvermt.TypeA,
		Answer: "127.0.0.1",
	})

	testAnswerB := HashAnswer(DNSAnswer{
		Type:   resolvermt.TypeCNAME,
		Answer: "bar",
	})

	testAnswerC := HashAnswer(DNSAnswer{
		Type:   resolvermt.TypeAAAA,
		Answer: "127.0.0.1",
	})

	tests := []struct {
		name  string
		haveA []AnswerHash
		haveB []AnswerHash
		want  bool
	}{
		{name: "empty records", haveA: []AnswerHash{}, haveB: []AnswerHash{}, want: false},
		{name: "empty B", haveA: []AnswerHash{testAnswerA}, haveB: []AnswerHash{}, want: false},
		{name: "not matching", haveA: []AnswerHash{testAnswerA}, haveB: []AnswerHash{testAnswerC}, want: false},
		{name: "matching", haveA: []AnswerHash{testAnswerA}, haveB: []AnswerHash{testAnswerA}, want: true},
		{name: "matching multiple", haveA: []AnswerHash{testAnswerA, testAnswerB, testAnswerC}, haveB: []AnswerHash{testAnswerC, testAnswerB}, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := answerMatch(test.haveA, test.haveB)

			assert.Equal(t, test.want, got)
		})
	}
}

func TestResolveWithCache(t *testing.T) {
	answer := DNSAnswer{Type: resolvermt.TypeA, Answer: "127.0.0.1"}
	answerHash := HashAnswer(answer)

	t.Run("no answers", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		got := task.resolveWithCache("test.com")
		assert.Empty(t, got)
	})

	t.Run("answer from cache", func(t *testing.T) {
		task, _ := newStubDetectionTask()
		task.ctx.dnsCache.Add("test.com", []DNSAnswer{answer})

		got := task.resolveWithCache("test.com")
		assert.ElementsMatch(t, []AnswerHash{answerHash}, got)
	})

	t.Run("answer from resolver", func(t *testing.T) {
		task, resolver := newStubDetectionTask()
		resolver.addAnswer("test.com", []DNSAnswer{answer})

		answers := task.resolveWithCache("test.com")
		assert.ElementsMatch(t, []AnswerHash{answerHash}, answers)
	})
}
