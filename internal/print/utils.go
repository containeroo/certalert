package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/kataras/tablewriter"
	"gopkg.in/yaml.v3"
)

// convertToYaml converts the output to yaml
func convertToYaml(output interface{}) (string, error) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(output); err != nil {
		return "", err
	}

	return b.String(), nil
}

// convertToJson converts the output to json
func convertToJson(output interface{}) (string, error) {
	var b bytes.Buffer
	jsonEncoder := json.NewEncoder(&b)
	jsonEncoder.SetIndent("", "  ") // set indent to 2 spaces
	if err := jsonEncoder.Encode(output); err != nil {
		return "", err
	}

	return b.String(), nil
}

// ConvertToTable converts the provided data to a table format.
func convertToTable(data interface{}) (string, error) {
	var output bytes.Buffer
	table := tablewriter.NewWriter(&output)

	s := reflect.ValueOf(data)

	if s.Kind() != reflect.Slice {
		return "", fmt.Errorf("Expected input of type slice for tabular conversion but received %s", s.Kind())
	}

	if s.Len() == 0 {
		return "", fmt.Errorf("Empty slice provided")
	}

	firstItem := s.Index(0)
	numFields := firstItem.NumField()
	var headers []string

	for i := 0; i < numFields; i++ {
		headers = append(headers, firstItem.Type().Field(i).Tag.Get("json"))
	}
	table.SetHeader(headers)

	for i := 0; i < s.Len(); i++ {
		item := s.Index(i)
		var row []string
		for j := 0; j < numFields; j++ {
			field := item.Field(j)
			row = append(row, fmt.Sprintf("%v", field.Interface()))
		}
		table.Append(row)
	}

	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()

	return output.String(), nil
}
