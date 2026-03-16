package workers

import (
	"errors"
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomSleep() error {
	sleep := time.Duration(rng.Intn(5)) * time.Second
	time.Sleep(sleep)

	if rng.Intn(2) == 0 {
		return errors.New("simulated worker failure")
	}

	return nil
}

func SendEmail(to, subject, body string) error {
	// Simulate sending an email
	return randomSleep()
}

func FetchURL(url string, timeout int) (string, error) {
	// Simulate fetching a URL
	return "fetched content", randomSleep()
}

func GenerateReport(userID, from, to string) (string, error) {
	// Simulate generating a report
	return "report content", randomSleep()
}
