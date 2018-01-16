package spectre

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

var tests map[string]*Test

// Init initially creates the list of Spectre Tests
func Init(testDirectory string) {
	tests = make(map[string]*Test)

	files, err := ioutil.ReadDir(testDirectory)
	if err != nil {
		fmt.Printf("Spectre Error: Unable to read test directory. %v\n", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fmt.Printf("Spectre: Found test definition file %v\n", file.Name())

			contents, err := ioutil.ReadFile(testDirectory + "/" + file.Name())
			if err != nil {
				fmt.Printf("Spectre Error: Unable to read definition file: %v", err)
				continue
			}

			var fileTest Test
			err = json.Unmarshal(contents, &fileTest)
			if err != nil {
				fmt.Printf("Spectre Error: Unable to parse definition file: %v", err)
				continue
			}

			tests[fileTest.ID] = &fileTest
		}
	}
}

// GetSpectreTestsForService returns the tests for the specified service
func GetSpectreTestsForService(serviceName string) []*Test {
	var results []*Test
	for _, test := range tests {
		if test.Service == serviceName {
			results = append(results, test)
		}
	}

	return results
}

// InvokeSpectreTest marks the test as invoked, reducing it's invocation count and removing it from the list if count is 0
func InvokeSpectreTest(ID string) error {
	test := getTestByID(ID)
	if test == nil {
		return errors.New("Spectre: Invalid Test")
	}

	test.InvocationCount--
	if test.InvocationCount < 1 {
		markTestComplete(ID)
	}

	return nil
}

func getTestByID(ID string) *Test {
	return tests[ID]
}

func markTestComplete(ID string) {
	for t, test := range tests {
		if test != nil && test.ID == ID {
			tests[t] = nil
		}
	}
}
