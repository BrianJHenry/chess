package chess

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Output to be used in tests
type MoveGenerationTestData struct {
	Description string
	Initial     State
	Results     []MoveGenerationResultTestData
}

type MoveGenerationResultTestData struct {
	Result      State
	Description string
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

func LoadTestData() ([]MoveGenerationTestData, error) {
	moveGenerationTestData := []MoveGenerationTestData{}
	err := filepath.WalkDir("../test_data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip non json files
		if filepath.Ext(path) != ".json" {
			return nil
		}

		testData, err := jsonToFormattedTestData(path)
		if err != nil {
			return err
		}

		moveGenerationTestData = append(moveGenerationTestData, testData...)
		return nil
	})

	return moveGenerationTestData, err
}

func jsonToFormattedTestData(filePath string) ([]MoveGenerationTestData, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return []MoveGenerationTestData{}, nil
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var fenTestCases FenTestCases
	json.Unmarshal(byteValue, &fenTestCases)

	var testData = make([]MoveGenerationTestData, 0, len(fenTestCases.TestCases))

	for _, testCase := range fenTestCases.TestCases {
		initialState, err := ConvertFenToState(testCase.Start.Fen)
		if err != nil {
			return []MoveGenerationTestData{}, err
		}

		var expected = make([]MoveGenerationResultTestData, 0, len(testCase.Expected))

		for _, expectedFenState := range testCase.Expected {
			expectedState, err := ConvertFenToState(expectedFenState.Fen)
			if err != nil {
				return []MoveGenerationTestData{}, err
			}

			expected = append(expected, MoveGenerationResultTestData{
				expectedState,
				expectedFenState.Move,
			})
		}

		testData = append(testData, MoveGenerationTestData{
			testCase.Start.Description,
			initialState,
			expected,
		})
	}

	return testData, nil
}
