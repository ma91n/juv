package example

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestModel(t *testing.T) {

	const payload = `
{"id":123,"Title":"test", "draft":true}
`
	var model Post
	err := json.Unmarshal([]byte(payload), &model)
	if err != nil {
		t.Fatal(err)
	}

	if model.ID != 123 {
		t.Fatal("expected 123 but actual is ", model.ID)
	}

	fmt.Println(model)

}
