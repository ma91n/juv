package example

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUnmarshallJSON(t *testing.T) {
	const payload = `{"id":123, "created":"2018-08-18T15:04:05Z", "title":"hoge", "body":"aaaaaaaa", "draft":true}`

	var p Post
	err := json.Unmarshal([]byte(payload), &p)
	if err != nil {
		t.Fatal(err)
	}

	if p.ID != 123 {
		t.Fatal("expected 123 but actual is ", p.ID)
	}
	if p.Title != "hoge" {
		t.Fatal("expected hoge but actual is ", p.Title)
	}
	if *p.Draft != true {
		t.Fatal("expected true but actual is ", p.Draft)
	}
	if p.Body != "aaaaaaaa" {
		t.Fatal("expected aaaaaaaa but actual is ", p.Body)
	}
	if p.Created.Unix() != 1534604645 {
		t.Fatal("expected 1566122078 but actual is ", p.Created.Unix())
	}
}

func TestValidateError(t *testing.T) {
	const payload = `{"id":-1, "title":"aaa", "draft":false}`

	var p Post
	err := json.Unmarshal([]byte(payload), &p)
	if err == nil {
		t.Fatal("test must occur validation error")
	}

	if !strings.Contains(err.Error(), "validation for 'ID' failed on the 'gt' tag") {
		t.Fatal("expected validation error but actual is", err.Error())
	}
}
