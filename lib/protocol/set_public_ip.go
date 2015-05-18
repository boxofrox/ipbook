package protocol

type SetPublicIpRequest struct {
	Name string `json:"name"`
}

func (r *SetPublicIpRequest) GetType() int {
	return TYPE_SET_PUBLIC_IP_REQUEST
}

func (r *SetPublicIpRequest) EncodeMessage() ([]byte, error) {
	return encodeMessage(r)
}

func (r *SetPublicIpRequest) ReadFrom(m *Message) error {
	return decode(m.Data, r)
}

type SetPublicIpResponse struct {
	Status Status `json:"status"`
	Reason string `json:"reason"`
}

func (r *SetPublicIpResponse) GetType() int {
	return TYPE_SET_IP_RESPONSE
}

func (r *SetPublicIpResponse) EncodeMessage() ([]byte, error) {
	return encodeMessage(r)
}

func (r *SetPublicIpResponse) ReadFrom(m *Message) error {
	return decode(m.Data, r)
}
