package chess

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Output to be used in tests
type MoveGenerationTestData struct {
	Initial State
	Results []State
}

// Data stuctures to hold json
type FenTestCases struct {
	Description string        `json:"description"`
	TestCases   []FenTestCase `json:"testCases"`
}

type FenTestCase struct {
	Start    FenStartState      `json:"start"`
	Expected []FenExpectedState `json:"expected"`
}

type FenStartState struct {
	Description string `json:"description"`
	Fen         string `json:"fen"`
}

type FenExpectedState struct {
	Move string `json:"move"`
	Fen  string `json:"fen"`
}

func JsonToFormattedTestData(filePath string) ([]MoveGenerationTestData, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return []MoveGenerationTestData{}, nil
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var fenTestCases FenTestCases
	json.Unmarshal(byteValue, &fenTestCases)

	var testData = make([]MoveGenerationTestData, 0, len(fenTestCases.TestCases))

	for _, testCase := range fenTestCases.TestCases {
		initialState, err := ConvertFenToState(testCase.Start.Fen)
		if err != nil {
			return []MoveGenerationTestData{}, err
		}

		var expected = make([]State, 0, len(testCase.Expected))

		for _, expectedFenState := range testCase.Expected {
			expectedState, err := ConvertFenToState(expectedFenState.Fen)
			if err != nil {
				return []MoveGenerationTestData{}, err
			}

			expected = append(expected, expectedState)
		}

		testData = append(testData, MoveGenerationTestData{
			initialState,
			expected,
		})
	}

	return testData, nil
}
