package id

import (
	"fmt"
	"testing"
)

func Test_IDService_maskAnyf(t *testing.T) {
	testCases := []struct {
		InputError  error
		InputFormat string
		InputArgs   []interface{}
		Expected    error
	}{
		{
			InputError:  nil,
			InputFormat: "",
			InputArgs:   []interface{}{},
			Expected:    nil,
		},
		{
			InputError:  fmt.Errorf("foo"),
			InputFormat: "bar",
			InputArgs:   []interface{}{},
			Expected:    nil,
		},
		{
			InputError:  fmt.Errorf("foo"),
			InputFormat: "bar %s",
			InputArgs:   []interface{}{"baz"},
			Expected:    fmt.Errorf("foo: bar baz"),
		},
	}

	for i, testCase := range testCases {
		var output error
		if len(testCase.InputArgs) == 0 {
			output = maskAnyf(testCase.InputError, testCase.InputFormat)
		} else {
			output = maskAnyf(testCase.InputError, testCase.InputFormat, testCase.InputArgs...)
		}

		if testCase.Expected != nil && output.Error() != testCase.Expected.Error() {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", output)
		}
	}
}
