package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplating(t *testing.T) {

	givenMemoryObject := make(map[string]interface{})
	givenMemoryObject["someKey"] = "Hello, World!"

	inputFile := []byte("{{.someKey}}")

	response, err := RunTemplate(inputFile, givenMemoryObject)

	if err != nil {
		t.Errorf("Test failed with error: %s", err)
	}
	assert.Equal(t, "Hello, World!", response)

}

func TestTemplatingWithDeepKeys(t *testing.T) {

	givenMemoryObject := make(map[string]interface{})
	givenMemoryObject["someKey"] = make(map[string]interface{})
	givenMemoryObject["someKey"].(map[string]interface{})["anotherKey"] = "Goodbye, World!"
	givenMemoryObject["someKey"].(map[string]interface{})["someKey"] = "Hello, World!"

	inputFile := []byte("{{.someKey.anotherKey}}")

	response, err := RunTemplate(inputFile, givenMemoryObject)

	if err != nil {
		t.Errorf("Test failed with error: %s", err)
	}

	assert.Equal(t, "Goodbye, World!", response)

}
