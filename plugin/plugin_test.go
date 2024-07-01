// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"fmt"
	"testing"
)

func TestSetAuthParams(t *testing.T) {
	tests := []struct {
		cmdArgs []string
		args    Args
		output  []string
		err     error
	}{
		// Test case 1
		{
			cmdArgs: []string{"executable", "arg1", "arg2"},
			args:    Args{Username: "john", Password: "password123", APIKey: "", AccessToken: ""},
			output:  []string{"executable", "arg1", "arg2", "--user $PLUGIN_USERNAME", "--password $PLUGIN_PASSWORD"},
			err:     nil,
		},
		// Test case 2
		{
			cmdArgs: []string{"./app", "--flag"},
			args:    Args{Username: "", Password: "", APIKey: "secretkey", AccessToken: ""},
			output:  []string{"./app", "--flag", "--apikey $PLUGIN_API_KEY"},
			err:     nil,
		},
		// Test case 3
		{
			cmdArgs: []string{"script.sh", "-option"},
			args:    Args{Username: "", Password: "", APIKey: "", AccessToken: "token123"},
			output:  []string{"script.sh", "-option", "--access-token $PLUGIN_ACCESS_TOKEN"},
			err:     nil,
		},
		// Test case 4
		{
			cmdArgs: []string{"command", "arg1"},
			args:    Args{Username: "", Password: "", APIKey: "", AccessToken: ""},
			output:  nil,
			err:     fmt.Errorf("either username/password, api key or access token needs to be set"),
		},
		// Test case 5
		{
			cmdArgs: []string{"app", "-flag"},
			args:    Args{Username: "user", Password: "", APIKey: "apikey123", AccessToken: ""},
			output:  []string{"app", "-flag", "--apikey $PLUGIN_API_KEY"},
			err:     nil,
		},
	}

	for _, tc := range tests {
		result, err := setAuthParams(tc.cmdArgs, tc.args)
		if err != nil {
			if tc.err == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, Got: %v", tc.err, err)
			}
		} else {
			if len(result) != len(tc.output) {
				t.Errorf("Expected output length: %d, Got: %d", len(tc.output), len(result))
			}
			for j := range result {
				if result[j] != tc.output[j] {
					t.Errorf("Mismatch at index %d. Expected: %s, Got: %s", j, tc.output[j], result[j])
				}
			}
		}
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		inputURL string
		expected string
		err      error
	}{
		{
			inputURL: "https://artifactory.maryout.com/artifactory/test44",
			expected: "https://artifactory.maryout.com/artifactory/",
			err:      nil,
		},
		{
			inputURL: "https://artifactory.maryout.com/artifactory/test/newdir/",
			expected: "https://artifactory.maryout.com/artifactory/",
			err:      nil,
		},
		{
			inputURL: "https://opautomates.jfrog.io/artifactory/test55/",
			expected: "https://opautomates.jfrog.io/artifactory/",
			err:      nil,
		},
		{
			inputURL: "https://opautomates.jfrog.io/artifactory",
			expected: "https://opautomates.jfrog.io/artifactory/",
			err:      nil,
		},
		{
			inputURL: "https://example.com/notartifactory",
			expected: "",
			err:      fmt.Errorf("url does not contain '/artifactory': https://example.com/notartifactory"),
		},
		{
			inputURL: "invalid-url",
			expected: "",
			err:      fmt.Errorf("invalid URL: invalid-url"),
		},
	}

	for _, tc := range tests {
		result, err := sanitizeURL(tc.inputURL)
		if err != nil {
			if tc.err == nil {
				t.Errorf("Unexpected error: %v", err)
			} else if err.Error() != tc.err.Error() {
				t.Errorf("Expected error: %v, Got: %v", tc.err, err)
			}
		} else {
			if result != tc.expected {
				t.Errorf("For URL %s, Expected: %s, Got: %s", tc.inputURL, tc.expected, result)
			}
		}
	}
}
