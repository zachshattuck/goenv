package goenv

import (
	"os"
	"testing"
)

type envTestSet struct {
	data        []byte
	expectError bool
	expectedEnv map[string]string
}

var envTests = []envTestSet{
	{
		data:        []byte("SERVER_HOST=localhost\nSERVER_PORT=7777\n"),
		expectError: false,
		expectedEnv: map[string]string{
			"SERVER_HOST": "localhost",
			"SERVER_PORT": "7777",
		},
	},
	{
		data:        []byte("SERVER_HOST=localhost\r\nSERVER_PORT=7777\n"),
		expectError: false,
		expectedEnv: map[string]string{
			"SERVER_HOST": "localhost",
			"SERVER_PORT": "7777",
		},
	},
	{
		data:        []byte("SOMETHING=nice\r\n\r\n\r\nELSE=notnice"),
		expectError: false,
		expectedEnv: map[string]string{
			"SOMETHING": "nice",
			"ELSE":      "notnice",
		},
	},
	{
		data:        []byte("OOPS\r\n\r\n\r\nELSE=notnice\r\n"),
		expectError: true,
		expectedEnv: map[string]string{},
	},
}

func TestDeserAndSetEnvironment(t *testing.T) {
	for i, test := range envTests {
		err := deserAndSetEnvironment(test.data)

		if err == nil && test.expectError {
			t.Errorf("[test %d] expected error, got nil", i)
		}

		if err != nil && !test.expectError {
			t.Errorf("[test %d] error: %v", i, err)
		}

		for key, value := range test.expectedEnv {
			if os.Getenv(key) != value {
				t.Errorf("[test %d] [%s] got '%v', want '%v'", i, key, os.Getenv(key), value)
			}
		}
	}
}

type readUntilTestSet struct {
	data        []byte
	startIdx    int
	delimiter   byte
	expectError bool
	expected    string
}

var readUntilTests = []readUntilTestSet{
	{
		data:        []byte("This is a nice string and all, AND it has a newline.\n"),
		startIdx:    0,
		delimiter:   '\n',
		expectError: false,
		expected:    "This is a nice string and all, AND it has a newline.",
	},
	{
		data:        []byte("SERVER_HOST=localhost\nSERVER_PORT=7777\n"),
		startIdx:    0,
		delimiter:   '\n',
		expectError: false,
		expected:    "SERVER_HOST=localhost",
	},
	{
		data:        []byte("SERVER_HOST=localhost\nSERVER_PORT=7777\n"),
		startIdx:    0,
		delimiter:   '=',
		expectError: false,
		expected:    "SERVER_HOST",
	},
	{
		data:        []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit"),
		startIdx:    0,
		delimiter:   '^',
		expectError: true,
		expected:    "",
	},
}

func TestReadUntil(t *testing.T) {
	for i, test := range readUntilTests {
		got, err := readUntil(test.data[test.startIdx:], test.delimiter, false)

		if err == nil && test.expectError {
			t.Errorf("[test %d] expected error, got nil", i)
		}

		if err != nil && !test.expectError {
			t.Errorf("[test %d] error: %v", i, err)
		}

		if string(got) != test.expected {
			t.Errorf("[test %d] got '%v', want '%v'", i, string(got), test.expected)
		}
	}
}
