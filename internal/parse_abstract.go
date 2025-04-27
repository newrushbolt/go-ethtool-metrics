package internal

import (
	"fmt"
	"log"
	"log/slog"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func parseFloat32(input string) float64 {
	floatRegex := regexp.MustCompile(`[+\-]?\d+(\.\d+)?`)
	floatString := floatRegex.FindString(input)
	float, _ := strconv.ParseFloat(floatString, 64)
	return float
}

func parseBool(input string) bool {
	output := false
	// TODO: better string to bool conversion, probably yaml-based?
	switch input {
	case "Yes":
		output = true
	case "On":
		output = true
	case "yes":
		output = true
	case "on":
		output = true
	}
	return output
}

func parseSlice(input string) []string {
	output := []string{}
	// TODO: drop /n stuff
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

func ParseAbstractDataObject(data *map[string]string, obj interface{}, tagName string) {
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()
	for key, value := range *data {
		if objType.Kind() == reflect.Pointer {
			slog.Warn("Skipping field because it's pointer. Something is probably broken in this data structure", "fieldName", key)
			continue
		}
		for i := 0; i < objValue.NumField(); i++ {
			field := objValue.Field(i)
			fieldType := objType.Field(i)

			object_tag := fieldType.Tag.Get(tagName)
			var tags []string
			if strings.Contains(object_tag, ",") {
				splitted_tags := strings.Split(object_tag, ",")
				tags = splitted_tags
			} else {
				if object_tag == "" {
					tags = []string{fieldType.Name}
				} else {
					tags = []string{object_tag}
				}
			}
			for _, tag := range tags {
				if tag == key {
					switch field.Kind() {
					case reflect.String:
						field.SetString(value)
					case reflect.Float32:
						field.SetFloat(parseFloat32(value))
					case reflect.Float64:
						field.SetFloat(parseFloat32(value))
					case reflect.Uint64:
						uint_value, err := strconv.ParseUint(value, 10, 64)
						if err != nil {
							log.Printf("Got uint64 parse error for key <%s:%s>: %s", tagName, tag, err)
							// TODO: DEBUG this line
							uint_value = 0
						}
						field.SetUint(uint_value)
					case reflect.Bool:
						field.SetBool(parseBool(value))
					case reflect.Slice:
						field.Set(reflect.ValueOf(parseSlice(value)))
					}
					// TODO: DEBUG if no type matched
				}
			}
		}
	}
}

func ParseAbstractColonData(data string, prefix string, keepPrefix bool) (map[string]string, error) {
	parsed_data := make(map[string]string)

	lines := make(map[int]string)
	pre_lines := strings.Split(data, "\n")
	for i, line := range pre_lines {
		if line == "" {
			continue
		}
		if strings.Contains(line, ":") {
			lines[i] = line
		} else {
			// Search for previous line with ":"
			prev_valid_index := -1
			prev_index := i
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
		if strings.Contains(line, "BusAddress") {
			slog.Debug("")
		}
		splittedLineLength := len(splittedLine)
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
		parsed_data[key] = value
	}
	return parsed_data, nil
}
