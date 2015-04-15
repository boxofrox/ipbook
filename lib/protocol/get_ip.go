package protocol

type GetIpRequest struct {
	Name string `json:"name"`
}

func (r *GetIpRequest) GetType() int {
	return TYPE_GET_IP_REQUEST
}

func (r *GetIpRequest) EncodeMessage() ([]byte, error) {
	return encodeMessage(r)
}

func (r *GetIpRequest) ReadFrom(m *Message) error {
	return decode(m.Data, r)
}

type GetIpResponse struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

func (r *GetIpResponse) GetType() int {
	return TYPE_GET_IP_RESPONSE
}

func (r *GetIpResponse) EncodeMessage() ([]byte, error) {
	return encodeMessage(r)
}

func (r *GetIpResponse) ReadFrom(m *Message) error {
	return decode(m.Data, r)
}
