package vectors

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// LoadStateTransition loads the state-transition.json into a
// StateTransitionVector instance
func LoadStateTransition() (StateTransitionVector, error) {
	var vector StateTransitionVector

	jsonFile, err := os.Open("state-transition.json")
	if err != nil {
		return vector, err
	}
	defer func() { _ = jsonFile.Close() }()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return vector, err
	}

	err = json.Unmarshal(bytes, &vector)
	if err != nil {
		return vector, err
	}

	return vector, nil
}
