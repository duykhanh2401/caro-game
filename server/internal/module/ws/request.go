package ws

type Request struct {
	ID       string                 `json:"id"`
	ClientID string                 `json:"clientId"`
	Body     map[string]interface{} `json:"body"`
	Type     RequestType            `json:"type"`
}

type RequestType int

const (
	GET_ROOMS RequestType = iota
	CHANGE_USERNAME
	JOIN_ROOM
	LEFT_ROOM
	SEND_MESSAGE
	GET_OLD_MESSAGES
	CREATE_ROOM
	GUEST_READY
	MASTER_READY
	GAME_HANDLE_REQUEST
)

func (t RequestType) String() string {
	return []string{
		"GET_ROOMS",
		"CHANGE_USERNAME",
		"JOIN_ROOM",
		"LEFT_ROOM",
		"SEND_MESSAGE",
		"GET_OLD_MESSAGES",
		"CREATE_ROOM",
		"GUEST_READY",
		"MASTER_READY",
		"GAME_HANDLE_REQUEST",
	}[t]
}
