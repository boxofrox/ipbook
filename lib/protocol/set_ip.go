package protocol

type SetIpRequest struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

func (r *SetIpRequest) GetType() int {
	return TYPE_SET_IP_REQUEST
}

func (r *SetIpRequest) EncodeMessage() ([]byte, error) {
	return encodeMessage(r)
}

func (r *SetIpRequest) ReadFrom(m *Message) error {
	return decode(m.Data, r)
}

type SetIpResponse struct {
	Status Status `json:"status"`
	Reason string `json:"reason"`
}

func (r *SetIpResponse) GetType() int {
	return TYPE_SET_IP_RESPONSE
}

func (r *SetIpResponse) EncodeMessage() ([]byte, error) {
	return encodeMessage(r)
}

func (r *SetIpResponse) ReadFrom(m *Message) error {
	return decode(m.Data, r)
}

type Status string

const (
	STATUS_OK    Status = "ok"
	STATUS_ERROR        = "error"
)

func (r Status) String() string {
	return string(r)
}
