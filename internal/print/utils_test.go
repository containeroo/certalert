package print

import (
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

func TestConvertToYaml(t *testing.T) {
	result, err := convertToYaml(sampleData)
	assert.Nil(t, err)
	expectedResult := "- id: 1\n  name: John\n- id: 2\n  name: Doe\n"
	assert.Equal(t, expectedResult, result)
}

func TestConvertToJson(t *testing.T) {
	result, err := convertToJson(sampleData)
	assert.Nil(t, err)

	expectedResult := "[\n  {\n    \"id\": 1,\n    \"name\": \"John\"\n  },\n  {\n    \"id\": 2,\n    \"name\": \"Doe\"\n  }\n]\n"
	assert.Equal(t, expectedResult, result)
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
