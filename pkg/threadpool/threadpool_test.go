package threadpool_test

import (
	"strings"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/threadpool"
	"github.com/stretchr/testify/assert"
)

type charCountTask struct {
	value string
	count int

	wantCount int
}

func (t *charCountTask) Run() {
	t.count = len(t.value)
}

func TestThreadPool(t *testing.T) {
	emptyList := []charCountTask{}

	singleList := []charCountTask{
		{value: "test", wantCount: 4},
	}

	smallList := []charCountTask{
		{value: "hello", wantCount: 5},
		{value: "world", wantCount: 5},
		{value: "foo", wantCount: 3},
		{value: "bar", wantCount: 3},
		{value: "test", wantCount: 4},
	}

	var bigList []charCountTask = make([]charCountTask, 1000)
	for i := range bigList {
		bigList[i].value = strings.Repeat("a", i+1)
		bigList[i].wantCount = i + 1
	}

	tests := []struct {
		name            string
		haveThreadCount int
		haveQueueSize   int
		haveTasks       []charCountTask
		wantQueries     int
	}{
		{name: "single worker", haveThreadCount: 1, haveQueueSize: 10, haveTasks: smallList},
		{name: "multiple workers", haveThreadCount: 3, haveQueueSize: 10, haveTasks: smallList},
		{name: "single queue", haveThreadCount: 3, haveQueueSize: 1, haveTasks: smallList},
		{name: "big list", haveThreadCount: 5, haveQueueSize: 1000, haveTasks: bigList},
		{name: "no tasks", haveThreadCount: 5, haveQueueSize: 10, haveTasks: emptyList},
		{name: "single task", haveThreadCount: 5, haveQueueSize: 10, haveTasks: singleList},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool := threadpool.NewThreadPool(test.haveThreadCount, test.haveQueueSize)
			defer pool.Close()

			for i := range test.haveTasks {
				pool.Execute(&test.haveTasks[i])
			}

			pool.Wait()

			for _, task := range test.haveTasks {
				assert.Equal(t, task.wantCount, task.count)
			}

			gotTotal := pool.CurrentCount()

			assert.Equal(t, len(test.haveTasks), gotTotal)
		})
	}
}
