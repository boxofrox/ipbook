package protocol

import (
	"encoding/json"
)

type SetIpRequest struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

func (req *SetIpRequest) GetType() int {
	return TYPE_SET_IP_REQUEST
}

func (req *SetIpRequest) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*req)
	if nil != err {
		return nil, err
	}

	msg := Message{
		Type: TYPE_SET_IP_REQUEST,
		Data: data,
	}
	return json.Marshal(&msg)
}

type SetIpResponse struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}

func (resp *SetIpResponse) GetType() int {
	return TYPE_SET_IP_RESPONSE
}

func (resp *SetIpResponse) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*resp)
	if nil != err {
		return nil, err
	}

	msg := Message{
		Type: TYPE_SET_IP_RESPONSE,
		Data: data,
	}
	return json.Marshal(&msg)
}

type Status string

const (
	STATUS_OK    Status = "ok"
	STATUS_ERROR                = "error"
)

func (r Status) String() string {
	return string(r)
}
