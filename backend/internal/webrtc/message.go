package webrtc

type Message struct {
	Event    string `json:"event"`
	Data     string `json:"data"`
	TargetID string `json:"targetId"`
	SenderID string `json:"senderId"`
}
