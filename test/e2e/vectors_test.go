package e2e

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVector(t *testing.T) {
	var vector StateTransitionVector

	jsonFile, err := os.Open("state-transition.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &vector)
	if err != nil {
		panic(err)
	}

	mJson, err := json.Marshal(vector)
	if err != nil {
		panic(err)
	}

	assert.JSONEq(t, string(bytes), string(mJson))
}
