package protocol

type Status string

const (
	STATUS_OK    Status = "ok"
	STATUS_ERROR        = "error"
)

func (r Status) String() string {
	return string(r)
}
