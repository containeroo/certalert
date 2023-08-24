package print

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock data
type Sample struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var sampleData = []Sample{
	{1, "John"},
	{2, "Doe"},
}

// Custom type that will throw an error when trying to marshal it into YAML
type BadType struct{}

func (bt *BadType) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("this is a bad type")
}

func (bt *BadType) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("this is a bad type")
}

func TestConvertToYaml(t *testing.T) {
	t.Run("Converts to yaml", func(t *testing.T) {
		result, err := convertToYaml(sampleData)
		assert.Nil(t, err)
		expectedResult := "- id: 1\n  name: John\n- id: 2\n  name: Doe\n"
		assert.Equal(t, expectedResult, result)
	})

	// Test error in YAML encoding
	t.Run("error in conversion", func(t *testing.T) {
		input := &BadType{}
		_, err := convertToYaml(input)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})
}

func TestConvertToJson(t *testing.T) {
	// Test successful conversion
	t.Run("successful conversion", func(t *testing.T) {
		input := map[string]string{"key": "value"}
		expected := `{
			"key": "value"
		}`
		output, err := convertToJson(input)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Cleanup the output and expected strings to remove whitespace and newlines
		cleanupOutput := strings.Map(func(r rune) rune {
			if r == '\n' || r == '\t' || r == ' ' {
				return -1
			}
			return r
		}, output)

		cleanupExpected := strings.Map(func(r rune) rune {
			if r == '\n' || r == '\t' || r == ' ' {
				return -1
			}
			return r
		}, expected)

		if cleanupOutput != cleanupExpected {
			t.Fatalf("Expected %s, got %s", expected, output)
		}
	})

	// Test error in JSON encoding
	t.Run("error in conversion", func(t *testing.T) {
		input := &BadType{}
		_, err := convertToJson(input)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})
}

func TestConvertToTable(t *testing.T) {
	result, err := convertToTable(sampleData)
	assert.Nil(t, err)

	expectedResult := "  ID  NAME  \n  1   John  \n  2   Doe   \n"
	assert.Equal(t, expectedResult, result)
}

func TestConvertToTableNotSlice(t *testing.T) {
	_, err := convertToTable(Sample{1, "John"})
	assert.NotNil(t, err)
	assert.Equal(t, "Expected input of type slice for tabular conversion but received struct", err.Error())
}

func TestConvertToTableEmptySlice(t *testing.T) {
	_, err := convertToTable([]Sample{})
	assert.NotNil(t, err)
	assert.Equal(t, "Empty slice provided", err.Error())
}
