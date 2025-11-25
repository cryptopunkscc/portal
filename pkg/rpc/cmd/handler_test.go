package cmd

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Json(t *testing.T) {
	expected := Handler{
		Func: "Func",
		Name: "Name",
	}
	j, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	actual := Handler{}
	err = json.Unmarshal(j, &actual)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, actual)
}
