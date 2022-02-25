package elizeutils

import "testing"

func TestSplitter(t *testing.T) {
	type testCase struct {
		input  string
		sep    string
		index  int
		output string
	}
	tests := []testCase{
		{"0:6:0", ":", 1, "6"},
		{"0:6:0", ":", 10, ""},
		{"0:6:0", "/", 0, "0:6:0"},
	}

	for _, tc := range tests {
		got := Splitter(tc.input, tc.sep, tc.index)
		if got != tc.output {
			t.Errorf("Expected %s and got %s", tc.output, got)
		}
	}

}
