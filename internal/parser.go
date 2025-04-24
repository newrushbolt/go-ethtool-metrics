package internal

import (
	"fmt"
	"log/slog"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func parseFloat64(input string) float64 {
	floatRegex := regexp.MustCompile(`[+\-]?\d+(\.\d+)?`)
	floatString := floatRegex.FindString(input)
	float, err := strconv.ParseFloat(floatString, 64)
	if err != nil {
		// TODO: consistent logging levels
		slog.Warn("Cannot parse string to float", "err", err, "string", input)
	}
	return float
}

func parseBool(input string) bool {
	output := false
	// TODO: better string to bool conversion, probably yaml-based?
	switch strings.ToUpper(input) {
	case "YES":
		output = true
	case "ON":
		output = true
	case "TRUE":
		output = true
	}
	return output
}

func parseSlice(input string) []string {
	output := []string{}
	// TODO: drop /n stuff (WTF this means? probably should rewrite func alltogether, with tests)
	input_lines := strings.Split(input, "\n")
	for _, line := range input_lines {
		input_columns := strings.Split(line, " ")
		for _, column := range input_columns {
			if column != "" {
				clean_column := strings.TrimSpace(strings.TrimSuffix(column, ","))
				output = append(output, clean_column)
			}
		}
	}
	return output
}

// Possible types in data are:
// - float64 (should be used for all numeric types, regarding of signed\unsigned and postitive\negative)
// - string
// - bool
// - slice
// - struct
func ParseAbstractDataObject(data *map[string]string, obj any, tagName string) {
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()
	for key, value := range *data {
		if objType.Kind() == reflect.Pointer {
			// TODO: consistent logging levels
			slog.Warn("Skipping field because it's pointer. Something is probably broken in this data structure", "fieldName", key)
			continue
		}
		for fieldIndex := range objValue.NumField() {
			fieldObj := objValue.Field(fieldIndex)
			fieldObjType := objType.Field(fieldIndex)

			object_tag := fieldObjType.Tag.Get(tagName)
			var tags []string
			if strings.Contains(object_tag, ",") {
				splitted_tags := strings.Split(object_tag, ",")
				tags = splitted_tags
			} else {
				if object_tag == "" {
					tags = []string{fieldObjType.Name}
				} else {
					tags = []string{object_tag}
				}
			}
			for _, tag := range tags {
				if tag == key {
					switch fieldObj.Kind() {
					case reflect.String:
						fieldObj.SetString(value)
					// Direct float values
					case reflect.Float64:
						fieldObj.SetFloat(parseFloat64(value))
					// Indirect float values, using pointer, and allowing for Nan-values
					case reflect.Pointer:
						realFieldType := fieldObjType.Type.Elem()
						if realFieldType.Kind() != reflect.Float64 {
							// TODO: consistent logging levels
							slog.Warn("Skipping field since only pointers to type float64 are supported", "field", fieldObjType.Name, "type", realFieldType.Kind())
							continue
						}
						floatObj := parseFloat64(value)
						fieldObj.Set(reflect.ValueOf(&floatObj))
					case reflect.Bool:
						fieldObj.SetBool(parseBool(value))
					case reflect.Slice:
						fieldObj.Set(reflect.ValueOf(parseSlice(value)))
					default:
						// TODO: consistent logging levels
						slog.Warn("Got unsupported data type in field", "field", key, "type", fieldObj.Kind())
					}
				}
			}
		}
	}
}

func ParseAbstractColonData(data string, prefix string, keepPrefix bool) map[string]string {
	parsed_data := make(map[string]string)

	lines := make(map[int]string)
	rawLines := strings.Split(data, "\n")
	for lineIndex, line := range rawLines {
		if line == "" {
			continue
		}
		if strings.Contains(line, ":") {
			lines[lineIndex] = line
		} else {
			// Search for previous line with ":"
			prev_valid_index := -1
			prev_index := lineIndex
			for prev_index > -1 {
				if _, ok := lines[prev_index-1]; ok {
					prev_valid_index = prev_index - 1
					break
				} else {
					prev_index -= 1
				}
			}
			if prev_valid_index != -1 {
				lines[prev_valid_index] = fmt.Sprintf("%s %s", lines[prev_valid_index], strings.TrimSpace(line))
			}
		}
	}
	for _, line := range lines {
		separatorIndex := strings.LastIndex(line, ": ")
		if separatorIndex < 0 {
			continue
		}
		splittedLine := []string{
			line[:separatorIndex],
			line[separatorIndex+1:],
		}
		splittedLineLength := len(splittedLine)
		// TODO: better logic for splitting lines
		if !(splittedLineLength == 2 || splittedLineLength == 3) {
			slog.Debug("Splitted line has invalid ammound of parts", "line", line, "splitted_line", splittedLine)
			continue
		}
		key := strings.TrimSpace(splittedLine[0])
		value := strings.TrimSpace(splittedLine[1])
		if prefix != "" {
			if strings.HasPrefix(key, prefix) {
				if !keepPrefix {
					key = strings.TrimSpace(strings.Replace(key, prefix, "", 1))
				}
			} else {
				continue
			}
		}
		// TODO: log elses?
		parsed_data[key] = value
	}
	return parsed_data
}
