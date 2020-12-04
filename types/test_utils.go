package types

import (
	"reflect"
	"testing"
)

func DeepEqual(t *testing.T, got interface{}, expected interface{}) {
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got = %v, want %v", got, expected)
	}

	gotVal := reflect.ValueOf(got).Elem()
	expectedVal := reflect.ValueOf(expected).Elem()
	for i := 0; i < gotVal.NumField(); i++ {
		valueField := gotVal.Field(i)
		typeField := gotVal.Type().Field(i)
		fieldName := typeField.Name

		gotValue := valueField.Interface()
		expectedValue := expectedVal.FieldByName(fieldName).Interface()
		if !reflect.DeepEqual(gotValue, expectedValue) {
			t.Errorf("%v: (%v/%v) expected: (%v/%v)", fieldName, gotValue, reflect.TypeOf(gotValue), expectedValue, reflect.TypeOf(expectedValue))
		}
	}
}
