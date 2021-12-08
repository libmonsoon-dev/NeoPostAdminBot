package repost

type Config struct {
	Sources     []string `json:"sources"`
	Destination string   `json:"destination"`

	DisableNotification bool `json:"disable_notification"`
	FromBackground      bool `json:"from_background"`
	SendCopy            bool `json:"send_copy"`
	RemoveCaption       bool `json:"remove_caption"`
	ReForward           bool `json:"re_forward"`
}
