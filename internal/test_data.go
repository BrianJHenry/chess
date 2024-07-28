package chess

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Output to be used in tests
type moveGenerationTestData struct {
	Description string
	Initial     State
	Results     []moveGenerationResultTestData
}

type moveGenerationResultTestData struct {
	Result State
	Move   string
}

// Data stuctures to hold json
type fenTestCases struct {
	Description string        `json:"description"`
	TestCases   []fenTestCase `json:"testCases"`
}

type fenTestCase struct {
	Start    fenStartState      `json:"start"`
	Expected []fenExpectedState `json:"expected"`
}

type fenStartState struct {
	Description string `json:"description"`
	Fen         string `json:"fen"`
}

type fenExpectedState struct {
	Move string `json:"move"`
	Fen  string `json:"fen"`
}

func loadTestData() ([]moveGenerationTestData, error) {
	moveGenerationTestData := []moveGenerationTestData{}
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

func jsonToFormattedTestData(filePath string) ([]moveGenerationTestData, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return []moveGenerationTestData{}, nil
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var fenTestCases fenTestCases
	json.Unmarshal(byteValue, &fenTestCases)

	var testData = make([]moveGenerationTestData, 0, len(fenTestCases.TestCases))

	for _, testCase := range fenTestCases.TestCases {
		initialState, err := fenToState(testCase.Start.Fen)
		if err != nil {
			return []moveGenerationTestData{}, err
		}

		var expected = make([]moveGenerationResultTestData, 0, len(testCase.Expected))

		for _, expectedFenState := range testCase.Expected {
			expectedState, err := fenToState(expectedFenState.Fen)
			if err != nil {
				return []moveGenerationTestData{}, err
			}

			expected = append(expected, moveGenerationResultTestData{
				expectedState,
				expectedFenState.Move,
			})
		}

		testData = append(testData, moveGenerationTestData{
			testCase.Start.Description,
			initialState,
			expected,
		})
	}

	return testData, nil
}
