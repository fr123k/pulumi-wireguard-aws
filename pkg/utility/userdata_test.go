package utility

import (
    "strings"
    "testing"
)

func assert(t *testing.T, err error) {
    if err != nil {
        t.Error(err.Error())
    }
}

func assertError(t *testing.T, err error, want string) {
    if err == nil {
        t.Errorf("Expected error was nil, got: nil, want: %s", want)
    }
    if !strings.Contains(err.Error(), want) {
        t.Errorf("The error message is wrong, got: %s, want: %s.", err.Error(), want)
    }
}

// TestReadFileWithNonExistingFile test the ReadFile method to read an non existing file expect error
func TestReadFileWithNonExistingFile(t *testing.T) {
    _, err := ReadFile("nonexistingfile.txt")

    assertError(t, err, "open nonexistingfile.txt: no such file or directory")
}

// TestReadFileFromMemory
func TestReadFileFromMemory(t *testing.T) {
    inMemReader := InMemoryFileReader{Str: "Memory Content"}
    content, err := NewInMemoryUtil(inMemReader).ReadFile("dsadasd")

    assert(t, err)

    if *content != inMemReader.Str {
        t.Errorf("The content is wrong, got: %s, want: %s.", *content, inMemReader.Str)
    }
}
