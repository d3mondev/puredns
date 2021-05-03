package sponsors

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShow_OK(t *testing.T) {
	// Restore stdout on exit
	stdOut := os.Stdout
	defer func() { os.Stdout = stdOut }()

	// Override stdout
	r, w, err := os.Pipe()
	require.Nil(t, err)
	os.Stdout = w

	// Stub http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "test")
	}))
	defer ts.Close()

	service := NewService()
	err = service.Show(ts.URL)
	assert.Nil(t, err)

	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	require.Nil(t, err)

	assert.Equal(t, "test\n", string(buf[:n]))
}

func TestShow_HTTPError(t *testing.T) {
	service := NewService()
	err := service.Show("")

	assert.NotNil(t, err)
}
