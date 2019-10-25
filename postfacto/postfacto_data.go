package postfacto

import "encoding/json"

func UnmarshalRetro(data []byte) (Retro, error) {
	var r Retro
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Retro) Marshal() ([]byte, error) {
	return json.Marshal(r)
}


type RetroItem struct {
	Description string   `json:"description"`
	Category    Category `json:"category"`
}

type passwordPayload struct {
	Retro retroPassword `json:"retro"`
}

type retroPassword struct {
	Password string `json:"password"`
}

type tokenReply struct {
	Token string `json:"token"`
}


type Retro struct {
	Retro RetroClass `json:"retro"`
}

type RetroClass struct {
	ID                int64         `json:"id"`
	Slug              string        `json:"slug"`
	Name              string        `json:"name"`
	IsPrivate         bool          `json:"is_private"`
	VideoLink         string        `json:"video_link"`
	CreatedAt         string        `json:"created_at"`
	HighlightedItemID interface{}   `json:"highlighted_item_id"`
	RetroItemEndTime  string        `json:"retro_item_end_time"`
	SendArchiveEmail  bool          `json:"send_archive_email"`
	Items             []interface{} `json:"items"`
	ActionItems       []ActionItem  `json:"action_items"`
	Archives          []Archive     `json:"archives"`
}

type ActionItem struct {
	ID          int64       `json:"id"`
	Description string      `json:"description"`
	Done        bool        `json:"done"`
	CreatedAt   string      `json:"created_at"`
	ArchivedAt  interface{} `json:"archived_at"`
}

type Archive struct {
	ID int64 `json:"id"`
}
