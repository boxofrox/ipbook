package protocol

import (
	"encoding/json"
)

type GetIpRequest struct {
	Name string `json:"name"`
}

func (req *GetIpRequest) GetType() int {
	return TYPE_GET_IP_REQUEST
}

func (req *GetIpRequest) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*req)
	if nil != err {
		return nil, err
	}

	msg := Message{
		Type: TYPE_GET_IP_REQUEST,
		Data: json.RawMessage(data),
	}

	return json.Marshal(&msg)
}

type GetIpResponse struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

func (resp *GetIpResponse) GetType() int {
	return TYPE_GET_IP_RESPONSE
}

func (resp *GetIpResponse) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*resp)
	if nil != err {
		return nil, err
	}

	msg := Message{
		Type: TYPE_GET_IP_RESPONSE,
		Data: data,
	}
	return json.Marshal(&msg)
}
