package testutil

import (
	"reflect"
	"strconv"
	"testing"
)

func TDeepEqual(t *testing.T, prompt string, got interface{}, expected interface{}) {
	valueOfGot := reflect.ValueOf(got)
	valueOfExpect := reflect.ValueOf(expected)

	gotKind := valueOfGot.Kind()
	expectedKind := valueOfExpect.Kind()

	if gotKind != expectedKind {
		t.Errorf("%v: kind: %v expected: %v", prompt, gotKind, expected)
		return
	}

	if gotKind == reflect.Slice || gotKind == reflect.Array {
		if gotKind == reflect.Slice {
			if valueOfGot.IsNil() && !valueOfExpect.IsNil() {
				t.Errorf("%v: is nil: %v expected: %v", prompt, got, expected)
				return
			}

			if !valueOfGot.IsNil() && valueOfExpect.IsNil() {
				t.Errorf("%v: not nil: %v expected: %v", prompt, got, expected)
				return
			}
		}

		lenGot := valueOfGot.Len()
		lenExpect := valueOfExpect.Len()
		if lenGot != lenExpect {
			t.Errorf("%v: kind: %v len: %v expected: %v", prompt, gotKind, lenGot, lenExpect)
		}
		if lenGot == 0 {
			return
		}

		zerothVal := valueOfGot.Index(0)
		zerothKind := zerothVal.Kind()
		if zerothKind != reflect.Ptr && zerothKind != reflect.Array && zerothKind != reflect.Slice {
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("%v: %v expected: %v", prompt, got, expected)
			}
			return
		}

		for i := 0; i < lenGot; i++ {
			gotVal := valueOfGot.Index(i)
			gotValue := gotVal.Interface()

			if i >= lenExpect {
				t.Errorf("%v (%v/%v): %v", prompt, i, lenGot, gotValue)
				continue
			}

			expectedVal := valueOfExpect.Index(i)

			tDeepEqualValue(t, prompt+":"+strconv.Itoa(i), gotVal, expectedVal)
		}
		return
	}

	if gotKind != reflect.Ptr && gotKind != reflect.Interface {
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("%v: (%v/%v) expected: (%v/%v)", prompt, got, gotKind, expected, expectedKind)
		}
		return
	}

	if valueOfGot.IsNil() && !valueOfExpect.IsNil() {
		t.Errorf("%v: is nil: %v expected: %v", prompt, got, expected)
		return
	}

	if !valueOfGot.IsNil() && valueOfExpect.IsNil() {
		t.Errorf("%v: (ptr) not nil: %v expected: %v", prompt, got, expected)
		return
	}

	if valueOfGot.IsNil() || valueOfExpect.IsNil() {
		return
	}

	gotVal := valueOfGot.Elem()
	expectedVal := valueOfExpect.Elem()

	gotKind = gotVal.Kind()
	expectedKind = expectedVal.Kind()

	if gotKind != reflect.Struct || expectedKind != reflect.Struct {
		tDeepEqualValue(t, prompt, gotVal, expectedVal)
		return
	}

	for i := 0; i < gotVal.NumField(); i++ {
		valueField := gotVal.Field(i)
		typeField := gotVal.Type().Field(i)
		fieldName := typeField.Name

		if fieldName[0] < 'A' || fieldName[0] > 'Z' {
			continue
		}

		expectedField := expectedVal.FieldByName(fieldName)

		tDeepEqualValue(t, prompt+":"+fieldName, valueField, expectedField)
	}
}

func tDeepEqualValue(t *testing.T, fieldName string, got reflect.Value, expected reflect.Value) {
	gotValue := got.Interface()
	expectedValue := expected.Interface()

	gotKind := got.Kind()
	expectedKind := expected.Kind()

	if gotKind != expectedKind {
		t.Errorf("%v: kind: %v expected: %v", fieldName, gotKind, expectedKind)
		return
	}

	if gotKind == reflect.Slice || gotKind == reflect.Array {
		TDeepEqual(t, fieldName, gotValue, expectedValue)
		return
	}

	if gotKind == reflect.Ptr {
		TDeepEqual(t, fieldName, gotValue, expectedValue)
		return
	}

	if gotKind == reflect.Interface {
		TDeepEqual(t, fieldName, gotValue, expectedValue)
		return
	}

	if !reflect.DeepEqual(gotValue, expectedValue) {
		t.Errorf("%v: %v expected: %v", fieldName, gotValue, expectedValue)
	}
}
