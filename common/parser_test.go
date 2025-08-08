package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFloat64(t *testing.T) {
	floatString := "3.14159265359"
	float := parseFloat64(floatString)
	// float64 in Go seems to be precise enough for comparing with Equal()
	assert.Equal(t, 3.14159265359, float)

	floatString = "3"
	float = parseFloat64(floatString)
	assert.Equal(t, 3.0, float)

	floatString = "+3"
	float = parseFloat64(floatString)
	assert.Equal(t, 3.0, float)

	floatString = "-3"
	float = parseFloat64(floatString)
	assert.Equal(t, -3.0, float)

	floatString = "some text 3.14 wont harm"
	float = parseFloat64(floatString)
	assert.Equal(t, 3.14, float)

	floatString = "nice text wont float"
	float = parseFloat64(floatString)
	assert.Equal(t, 0.0, float)

	floatString = "3.14 two floats / 3.15 one string"
	float = parseFloat64(floatString)
	assert.Equal(t, 3.14, float)
}

func TestParseBoolTrue(t *testing.T) {
	boolIsh := "On"
	boolVal := parseBool(boolIsh)
	assert.True(t, boolVal)

	boolIsh = "ON"
	boolVal = parseBool(boolIsh)
	assert.True(t, boolVal)

	boolIsh = " on "
	boolVal = parseBool(boolIsh)
	assert.True(t, boolVal)

	boolIsh = "TRUE"
	boolVal = parseBool(boolIsh)
	assert.True(t, boolVal)

	boolIsh = "yes"
	boolVal = parseBool(boolIsh)
	assert.True(t, boolVal)
}

func TestParseBoolFalse(t *testing.T) {
	boolIsh := "no"
	boolVal := parseBool(boolIsh)
	assert.False(t, boolVal)

	boolIsh = "NO"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)

	boolIsh = "off"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)

	boolIsh = "oFF"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)

	boolIsh = "false"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)

	// Shouldn't be parsed
	boolIsh = " ON OFF"
	boolVal = parseBool(boolIsh)
	assert.False(t, boolVal)
}

func TestParseSliceMulti(t *testing.T) {
	expectedSlice := []string{"mode1", "mode2", "mode3"}
	sliceString := "mode1, mode2, mode3"
	parsedSlice := parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "mode1 mode2 mode3"
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "mode1,mode2,mode3"
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "mode1\nmode2\nmode3"
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)
}

func TestParseSliceEmpty(t *testing.T) {
	expectedSlice := []string{}
	sliceString := ""
	parsedSlice := parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "  "
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "  \n \t ,  \t"
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = "Not reported"
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)

	sliceString = " Not reported "
	parsedSlice = parseSlice(sliceString)
	assert.Equal(t, expectedSlice, parsedSlice)
}

func TestParseString(t *testing.T) {
	assert.Equal(t, "hello", parseString(" hello "))
	assert.Equal(t, "", parseString("Not reported"))
	assert.Equal(t, "", parseString(" None "))
	assert.Equal(t, "some value", parseString("some value"))
	assert.Equal(t, "", parseString("  "))
}

func TestParseAbstractDataObject(t *testing.T) {
	type TestStruct struct {
		Str         string   `testtag:"tag1"`
		Flt         float64  `testtag:"tag2"`
		Bool        bool     `testtag:"tag3"`
		Slice       []string `testtag:"tag4"`
		PtrFlt      *float64 `testtag:"tag5"`
		PtrStr      *string  `testtag:"tag6"`
		PtrPtr      **string `testtag:"tag10"`
		NoTag       string
		MultiTag    string            `testtag:"tag7,tag8"`
		Unsupported map[string]string `testtag:"tag9"`
	}

	obj := &TestStruct{}
	data := map[string]string{
		"tag1":  "hello",
		"tag2":  "123.45",
		"tag3":  "true",
		"tag4":  "a, b, c",
		"tag5":  "99.9",
		"tag6":  "should not set",
		"NoTag": "fieldname",
		"tag7":  "multi",
		"tag9":  "unsupported",
		"tag10": "WTF",
	}
	ParseAbstractDataObject(&data, obj, "testtag")

	assert.Equal(t, "hello", obj.Str)
	assert.Equal(t, 123.45, obj.Flt)
	assert.True(t, obj.Bool)
	assert.Equal(t, 3, len(obj.Slice))
	assert.NotNil(t, obj.PtrFlt)
	assert.Equal(t, 99.9, *obj.PtrFlt)
	assert.Nil(t, obj.PtrStr)
	assert.Equal(t, "fieldname", obj.NoTag)
	assert.Equal(t, "multi", obj.MultiTag)
	assert.Nil(t, obj.Unsupported)
	assert.Nil(t, obj.PtrPtr)
}

func TestParseAbstractDataObject_StructPointerType(t *testing.T) {
	type PtrStruct struct {
		Field string
	}
	obj := &PtrStruct{}
	data := map[string]string{"Field": "value"}
	// Pass pointer to struct type, triggers objType.Kind() == reflect.Pointer branch
	ParseAbstractDataObject(&data, &obj, "testtag")
	// Should skip setting any fields, obj remains unchanged
	assert.Equal(t, "", obj.Field)
}

func TestParseAbstractColonData_InvalidSplit(t *testing.T) {
	expectedResult := map[string]string{"foo: bar: baz": "qux"}
	input := "foo: bar: baz: qux"
	result := ParseAbstractColonData(input, "", false)
	assert.Equal(t, expectedResult, result)
}
