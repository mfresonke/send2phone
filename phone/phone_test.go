package phone

import (
	"os"
	"testing"

	"github.com/mfresonke/ngrokker"
)

const (
	testingPort              = 7070
	testingVerbose           = false
	testingValidNumber       = "+14071111111"
	testingAcceptedNGROKTOS  = true
	testingAcceptedTwilioTOS = true
)

type testTunnel struct{}

func (tt testTunnel) Open(port int) ([]ngrokker.Endpoint, error) {
	return nil, nil
}

func (tt testTunnel) Close() error {
	return nil
}

func testingSender() *Sender {
	return NewSenderTunnel(
		testTunnel{},
		TwilioConfig{
			SID:       "SID",
			AuthToken: "Auth",
			SenderNum: "SomeValidNum",
		},
		testingPort,
		testingVerbose,
	)
}

func expectError(t *testing.T, err, expectedErr error, optionalMsg string) {
	if err == nil {
		t.Fatal("No error returned")
	}
	if err != expectedErr {
		if optionalMsg == "" {
			t.Error("An error was returned, but it was of the wrong type.")
		} else {
			t.Error(optionalMsg)
		}
	}
}

func TestErrFileDoesNotExist(t *testing.T) {
	sender := testingSender()
	err := sender.SendFile(testingValidNumber, "/this/does/not/exist")
	expectError(t, err, ErrFileDoesNotExist, "Putting an invalid file did not properly return ErrFileDoesNotExist")
}

func TestErrFiletypeNotSupported(t *testing.T) {
	invalidFilePath := os.TempDir() + "/invalidFile.lol"
	_, err := os.Create(invalidFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidFilePath)
	sender := testingSender()
	err = sender.SendFile(testingValidNumber, invalidFilePath)
	expectError(t, err, ErrFiletypeNotSupported, "Putting an invalid file did not properly return ErrFiletypeNotSupported")
}

func TestErrFileIsDirectory(t *testing.T) {
	sender := testingSender()
	err := sender.SendFile(testingValidNumber, os.TempDir())
	expectError(t, err, ErrFileIsDirectory, "ErrFileIsDirectory not properly returned.")
}
