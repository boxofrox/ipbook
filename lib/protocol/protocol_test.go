package protocol

import (
	"testing"
)

func Test_decode_WithGetIpRequestJsonMessage_GetIpRequestObject(t *testing.T) {
	expected := &GetIpRequest{"home"}

	encoding, _ := encode(expected)
	object, err := Decode(encoding)

	if nil != err {
		t.Fatalf("error while decoding. %q.\nJSON: %s", err.Error(), string(encoding))
	}

	if nil == object {
		t.Fatalf("expecting valid Messager type, got nil")
	}

	actual := object.(*GetIpRequest)

	if actual.GetType() != expected.GetType() {
		t.Errorf("expecting type %q, got %q", expected.GetType(), actual.GetType())
	}

	if actual.Name != expected.Name {
		t.Errorf("expecting Name %q, got %q", expected.Name, actual.Name)
	}
}

func Test_decode_WithGetIpResponseJsonMessage_GetIpResponseObject(t *testing.T) {
	expected := &GetIpResponse{"home", "0.0.0.0"}

	encoding, _ := encode(expected)
	object, err := Decode(encoding)

	if nil != err {
		t.Fatalf("error while decoding. %q.\nJSON: %s", err.Error(), string(encoding))
	}

	if nil == object {
		t.Fatalf("expecting valid Messager type, got nil")
	}

	actual := object.(*GetIpResponse)

	if actual.GetType() != expected.GetType() {
		t.Errorf("expecting type %q, got %q", expected.GetType(), actual.GetType())
	}

	if actual.Name != expected.Name {
		t.Errorf("expecting Name %q, got %q", expected.Name, actual.Name)
	}

	if actual.Ip != expected.Ip {
		t.Errorf("expecting IP %q, got %q", expected.Ip, actual.Ip)
	}
}

func Test_decode_WithSetIpRequestJsonMessage_SetIpRequestObject(t *testing.T) {
	expected := &SetIpRequest{"home", "0.0.0.0"}

	encoding, _ := encode(expected)
	object, err := Decode(encoding)

	if nil != err {
		t.Fatalf("error while decoding. %q.\nJSON: %s", err.Error(), string(encoding))
	}

	if nil == object {
		t.Fatalf("expecting valid Messager type, got nil")
	}

	actual := object.(*SetIpRequest)

	if actual.GetType() != expected.GetType() {
		t.Errorf("expecting type %q, got %q", expected.GetType(), actual.GetType())
	}

	if actual.Name != expected.Name {
		t.Errorf("expecting Name %q, got %q", expected.Name, actual.Name)
	}

	if actual.Ip != expected.Ip {
		t.Errorf("expecting IP %q, got %q", expected.Ip, actual.Ip)
	}
}

func Test_decode_WithSetIpResponseJsonMessage_SetIpResponse(t *testing.T) {
	expected := &SetIpResponse{STATUS_OK, ""}

	encoding, _ := encode(expected)
	object, err := Decode(encoding)

	if nil != err {
		t.Fatalf("error while decoding. %q.\nJSON: %s", err.Error(), string(encoding))
	}

	if nil == object {
		t.Fatalf("expecting valid Messager type, got nil")
	}

	actual := object.(*SetIpResponse)

	if actual.GetType() != expected.GetType() {
		t.Errorf("expecting type %q, got %q", expected.GetType(), actual.GetType())
	}

	if actual.Status != expected.Status {
		t.Errorf("expecting Status %q, got %q", expected.Status, actual.Status)
	}
}

func Test_encode_WithValidMessage_NoErrors(t *testing.T) {
	tests := []struct {
		input    messager
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
