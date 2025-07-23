package common

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

var (
	Logger *slog.Logger
)

func init() {
	loggerLever := GetLogLevel()
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: loggerLever}))
}

// Parses and returns the first float found in string
// Scientific notations are not supported
func parseFloat64(input string) float64 {
	floatRegex := regexp.MustCompile(`[+\-]?\d+(\.\d+)?`)
	floatString := floatRegex.FindString(input)

	float, err := strconv.ParseFloat(floatString, 64)
	if err != nil {
		// TODO: consistent logging levels
		Logger.Warn("Cannot parse string to float", "err", err, "string", input)
	}
	return float
}

// Parses and returns the bool, using YAML loose convertion
// Read more on what's valid bool in yaml here
// https://github.com/go-yaml/yaml/blob/v3.0.1/decode.go#L678
func parseBool(input string) bool {
	var result bool
	err := yaml.Unmarshal([]byte(input), &result)
	if err != nil {
		Logger.Warn("Cannot parse bool (~ish) string to actual bool", "err", err, "string", input)
	}
	return result
}

func parseSlice(input string) []string {
	splitFunc := func(c rune) bool {
		return unicode.IsSpace(c) || c == ','
	}
	output := strings.FieldsFunc(input, splitFunc)
	// // Shippet for cleaning excess items from slice
	// deleteFunc := func (s string) bool {
	// 	return true
	// }
	// slices.DeleteFunc()output

	return output
}

// Possible types in data are:
// - float64 (should be used for all numeric types, regarding of signed\unsigned and positive\negative)
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
			Logger.Warn("Skipping field because it's pointer. Something is probably broken in this data structure", "fieldName", key)
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
					// Indirect values, using pointer
					case reflect.Pointer:
						realFieldType := fieldObjType.Type.Elem()
						switch realFieldType.Kind() {
						// ATM looks like it only makes sense to nil-fy floats,
						// since empty strings in prometheus labels are treated as missing
						case reflect.Float64:
							floatObj := parseFloat64(value)
							fieldObj.Set(reflect.ValueOf(&floatObj))
						default:
							// TODO: consistent logging levels
							Logger.Warn("Skipping field since only pointers to type float64 are supported", "field", fieldObjType.Name, "type", realFieldType.Kind())
							continue
						}
					case reflect.Bool:
						fieldObj.SetBool(parseBool(value))
					case reflect.Slice:
						fieldObj.Set(reflect.ValueOf(parseSlice(value)))
					default:
						// TODO: consistent logging levels
						Logger.Warn("Got unsupported data type in field", "field", key, "type", fieldObj.Kind())
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
			Logger.Debug("Splitted line has invalid amount of parts", "line", line, "splitted_line", splittedLine)
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
