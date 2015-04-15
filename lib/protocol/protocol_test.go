package protocol

import (
	"bytes"
	"testing"
)

func Test_decode_WithSimpleMessage_GetMessageObject(t *testing.T) {
	expected := &Message{TYPE_ERROR_RESPONSE, []byte("{\"key\":\"value\"}")}

	encoding, err := encode(expected)

	if nil != err {
		t.Fatalf("error while encoding message. %s", err)
	}

	actual := &Message{}

	if err := decode(encoding, actual); nil != err {
		t.Fatalf("error while decoding message. %s", err)
	}

	if actual.Type != expected.Type {
		t.Errorf("expecting Type %v, got %v", expected.Type, actual.Type)
	}

	if !bytes.Equal(actual.Data, expected.Data) {
		t.Errorf("expecting Data %v, got %v", expected.Data, actual.Data)
	}
}

func Test_decode_WithGetIpRequestJsonMessage_GetIpRequestObject(t *testing.T) {
	expected := &GetIpRequest{"home"}
	message := &Message{}

	encoding, _ := expected.EncodeMessage()

	if err := decode(encoding, message); nil != err {
		t.Fatalf("error while decoding. %s", err.Error())
	}

	if message.Type != expected.GetType() {
		t.Fatalf("expecting type %q, got %q", expected.GetType(), message.Type)
	}

	actual := &GetIpRequest{}

	if err := actual.ReadFrom(message); nil != err {
		t.Fatalf("error while decoding. %s", err.Error())
	}

	if actual.Name != expected.Name {
		t.Errorf("expecting Name %q, got %q", expected.Name, actual.Name)
	}
}

func Test_decode_WithGetIpResponseJsonMessage_GetIpResponseObject(t *testing.T) {
	expected := &GetIpResponse{"home", "0.0.0.0"}
	message := &Message{}

	encoding, _ := expected.EncodeMessage()

	if err := decode(encoding, message); nil != err {
		t.Fatalf("error while decoding. %s", err.Error())
	}

	if message.Type != expected.GetType() {
		t.Fatalf("expecting type %q, got %q", expected.GetType(), message.Type)
	}

	actual := &GetIpResponse{}

	if err := actual.ReadFrom(message); nil != err {
		t.Fatalf("error while decoding. %s", err.Error())
	}

	if actual.Name != expected.Name {
		t.Errorf("expecting Name %q, got %q", expected.Name, actual.Name)
	}

	if actual.Ip != expected.Ip {
		t.Errorf("expecting Ip %q, got %q", expected.Ip, actual.Ip)
	}
}

func Test_decode_WithSetIpRequestJsonMessage_SetIpRequestObject(t *testing.T) {
	expected := &SetIpRequest{"home", "0.0.0.0"}
	message := &Message{}

	encoding, _ := expected.EncodeMessage()

	if err := decode(encoding, message); nil != err {
		t.Fatalf("error while decoding. %s", err.Error())
	}

	actual := &SetIpRequest{}

	if err := actual.ReadFrom(message); nil != err {
		t.Fatalf("error while decoding. %s", err.Error())
	}

	if actual.Name != expected.Name {
		t.Errorf("expecting Name %q, got %q", expected.Name, actual.Name)
	}

	if actual.Ip != expected.Ip {
		t.Errorf("expecting IP %q, got %q", expected.Ip, actual.Ip)
	}
}

func Test_decode_WithSetIpResponseJsonMessage_SetIpResponse(t *testing.T) {
	expected := &SetIpResponse{STATUS_OK, "reason"}
	message := &Message{}

	encoding, _ := expected.EncodeMessage()

	if err := decode(encoding, message); nil != err {
		t.Fatalf("error while decoding. %s", err.Error())
	}

	actual := &SetIpResponse{}

	if err := actual.ReadFrom(message); nil != err {
		t.Fatalf("error while decoding. %s", err.Error())
	}

	if actual.Reason != expected.Reason {
		t.Errorf("expecting Reason %q, got %q", expected.Reason, actual.Reason)
	}

	if actual.Status != expected.Status {
		t.Errorf("expecting Status %q, got %q", expected.Status, actual.Status)
	}
}

func Test_encode_WithValidMessage_NoErrors(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected error
	}{
		{&GetIpRequest{"home"}, nil},
		{&GetIpResponse{"home", "0.0.0.0"}, nil},
		{&SetIpRequest{"home", "0.0.0.0"}, nil},
		{&SetIpResponse{STATUS_OK, ""}, nil},
	}

	for i, test := range tests {
		_, actual := encode(test.input)
		if test.expected != actual {
			t.Errorf("%d: expected %q, got %q", i, test.expected, actual)
		}
	}
}
