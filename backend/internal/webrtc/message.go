package webrtc

type Message struct {
	Event    string `json:"event"`
	Data     any    `json:"data"`
	TargetID string `json:"targetId"`
	SenderID string `json:"senderId"`
}
