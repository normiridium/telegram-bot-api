package tgbotapi

import (
	"testing"
)

func assertLen(t *testing.T, params Params, l int) {
	actual := len(params)
	if actual != l {
		t.Fatalf("Incorrect number of params, expected %d but found %d\n", l, actual)
	}
}

func assertEq(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("Values did not match, a: %v, b: %v\n", a, b)
	}
}

func TestAddNonEmpty(t *testing.T) {
	params := make(Params)
	params.AddNonEmpty("value", "value")
	assertLen(t, params, 1)
	assertEq(t, params["value"], "value")
	params.AddNonEmpty("test", "")
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddNonZero(t *testing.T) {
	params := make(Params)
	params.AddNonZero("value", 1)
	assertLen(t, params, 1)
	assertEq(t, params["value"], "1")
	params.AddNonZero("test", 0)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddNonZero64(t *testing.T) {
	params := make(Params)
	params.AddNonZero64("value", 1)
	assertLen(t, params, 1)
	assertEq(t, params["value"], "1")
	params.AddNonZero64("test", 0)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddBool(t *testing.T) {
	params := make(Params)
	params.AddBool("value", true)
	assertLen(t, params, 1)
	assertEq(t, params["value"], "true")
	params.AddBool("test", false)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddNonZeroFloat(t *testing.T) {
	params := make(Params)
	params.AddNonZeroFloat("value", 1)
	assertLen(t, params, 1)
	assertEq(t, params["value"], "1.000000")
	params.AddNonZeroFloat("test", 0)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddInterface(t *testing.T) {
	params := make(Params)
	data := struct {
		Name string `json:"name"`
	}{
		Name: "test",
	}
	params.AddInterface("value", data)
	assertLen(t, params, 1)
	assertEq(t, params["value"], `{"name":"test"}`)
	params.AddInterface("test", nil)
	assertLen(t, params, 1)
	assertEq(t, params["test"], "")
}

func TestAddFirstValid(t *testing.T) {
	params := make(Params)
	params.AddFirstValid("value", 0, "", "test")
	assertLen(t, params, 1)
	assertEq(t, params["value"], "test")
	params.AddFirstValid("value2", 3, "test")
	assertLen(t, params, 2)
	assertEq(t, params["value2"], "3")
}

func TestRichMessageConfigParams(t *testing.T) {
	config := NewRichMessageMarkdown(42, "# Title\n\n- item")
	config.ReplyToMessageID = 7
	config.AllowSendingWithoutReply = true

	params, err := config.params()
	if err != nil {
		t.Fatalf("unexpected params error: %v", err)
	}
	assertEq(t, config.method(), "sendRichMessage")
	assertEq(t, params["chat_id"], "42")
	assertEq(t, params["rich_message"], `{"markdown":"# Title\n\n- item"}`)
	assertEq(t, params["reply_parameters"], `{"message_id":7,"allow_sending_without_reply":true}`)
	if params["reply_to_message_id"] != "" {
		t.Fatalf("sendRichMessage must use reply_parameters, got legacy reply_to_message_id=%q", params["reply_to_message_id"])
	}
}

func TestRichMessageConfigParamsWithUploadMedia(t *testing.T) {
	config := NewRichMessageHTMLWithMedia(42, `<img src="tg://photo?id=formula_1">`, []InputRichMessageMedia{
		NewInputRichMessageMediaPhoto("formula_1", FileBytes{Name: "formula.png", Bytes: []byte{1, 2, 3}}),
	})

	params, err := config.params()
	if err != nil {
		t.Fatalf("unexpected params error: %v", err)
	}
	assertEq(t, config.method(), "sendRichMessage")
	assertEq(t, params["chat_id"], "42")
	assertEq(t, params["rich_message"], `{"html":"\u003cimg src=\"tg://photo?id=formula_1\"\u003e","media":[{"id":"formula_1","media":{"type":"photo","media":"attach://file-0","caption_entities":null}}]}`)

	files := config.files()
	if len(files) != 1 {
		t.Fatalf("expected 1 upload file, got %d", len(files))
	}
	assertEq(t, files[0].Name, "file-0")
	if _, ok := files[0].Data.(FileBytes); !ok {
		t.Fatalf("expected FileBytes upload, got %T", files[0].Data)
	}
}
