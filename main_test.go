package main

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		title string
		file  string
		err   string
	}{
		{
			title: "valid config",
			file:  "./testdata/valid-config.json",
		},
		{
			title: "invalid config",
			file:  "./testdata/invalid-config.json",
			err:   "invalid character 'I' looking for beginning of value",
		},
		{
			title: "not existing file",
			file:  "./testdata/nothing-here.json",
			err:   "open ./testdata/nothing-here.json: no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var fc FileConfig
			err := loadFileConfig(tc.file, &fc)
			errMsg := fmt.Sprintf("%v", err)

			if err != nil && tc.err == "" {
				t.Errorf(`Expected to NOT fail, got "%v"`, errMsg)
			} else if err != nil && tc.err != errMsg {
				t.Errorf(`Expected to fail, got "%v" but want "%v"`, errMsg, tc.err)
			}
		})
	}
}
