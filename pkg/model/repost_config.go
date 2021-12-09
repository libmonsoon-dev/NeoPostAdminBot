package model

type RepostConfig struct {
	Source      string `json:"sources"`
	Destination string `json:"destination"`

	SourceId      int64 `json:"source_id"`
	DestinationId int64 `json:"destination_id"`

	DisableNotification bool `json:"disable_notification"`
	FromBackground      bool `json:"from_background"`
	SendCopy            bool `json:"send_copy"`
	RemoveCaption       bool `json:"remove_caption"`
	ReForward           bool `json:"re_forward"`
}
