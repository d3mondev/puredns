package sponsors

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Service is a Sponsors service.
type Service struct{}

// NewService creates a new Service.
func NewService() *Service {
	return &Service{}
}

// Show downloads the sponsors text file and sends it to the console.
func (s Service) Show(fileURL string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	fmt.Println()

	return err
}
