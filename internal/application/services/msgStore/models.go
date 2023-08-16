package msgStore

type MsgData struct {
	ID       int         `json:"id"`
	Type     string      `json:"type"` //private_supergroup
	Name     string      `json:"name"`
	Messages []MsgImport `json:"messages"`
}

// MsgImport
// "action": "pin_message" => + field "message_id": 522,
// "action": "topic_edit" => + field "new_title": "xx", + field "new_icon_emoji_id": 0,
// "media_type": "sticker", => + field "sticker_emoji": "😍",
type MsgImport struct {
	ID             int          `json:"id"`
	Type           string       `json:"type,omitempty"` //service message
	Date           string       `json:"date,omitempty"`
	DateUnixtime   string       `json:"date_unixtime,omitempty"`
	Edited         string       `json:"edited,omitempty"`
	EditedUnixtime string       `json:"edited_unixtime,omitempty"`
	From           string       `json:"from,omitempty"`
	FromID         string       `json:"from_id,omitempty"`
	Text           interface{}  `json:"text,omitempty"`
	TextEntities   []TextEntity `json:"text_entities,omitempty"`
	ReplyTo        int          `json:"reply_to_message_id,omitempty"`
	Actor          string       `json:"actor,omitempty"` //
	ActorID        string       `json:"actor_id,omitempty"`
	Action         string       `json:"action,omitempty"` //Action: join_group_by_link remove_members invite_members migrate_from_group pin_message topic_created topic_edit edit_group_title edit_group_photo
	MessageID      int          `json:"message_id,omitempty"`
	Inviter        string       `json:"inviter,omitempty"` //Group
	Members        []string     `json:"members,omitempty"`
	File           string       `json:"file,omitempty"`
	Photo          string       `json:"photo,omitempty"`
	Thumbnail      string       `json:"thumbnail,omitempty"`
	StickerEmoji   string       `json:"sticker_emoji,omitempty"` // for media_type sticker
	Height         string       `json:"height,omitempty"`
	Width          string       `json:"width,omitempty"`
	MediaType      string       `json:"media_type,omitempty"`       // video_file animation audio_file sticker sticker voice_message
	MimeType       string       `json:"mime_type,omitempty"`        // video/mp4
	DurationSec    int          `json:"duration_seconds,omitempty"` //for media
	Performer      string       `json:"performer,omitempty"`        // for audio_file
	Title          string       `json:"title,omitempty"`            // for audio_file
	ContactInf     ContactData  `json:"contact_information,omitempty"`
	LocationInf    LocationData `json:"location_information,omitempty"`         //location
	LocationTime   int          `json:"live_location_period_seconds,omitempty"` //location
	PlaceName      string       `json:"place_name,omitempty"`                   // for location
	Address        string       `json:"address,omitempty"`                      // for location
	Poll           PollData     `json:"poll,omitempty"`
	ViaBot         string       `json:"via_bot,omitempty"` // for bot msg
}

type PollData struct {
	Question    string     `json:"question,omitempty"`
	TotalVoters int        `json:"total_voters,omitempty"`
	Closed      bool       `json:"closed,omitempty"`
	Answers     PollAnswer `json:"answers,omitempty"`
}

type PollAnswer struct {
	Text   string `json:"text,omitempty"`
	Voters int    `json:"voters,omitempty"`
	Chosen bool   `json:"chosen,omitempty"`
}

type ContactData struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type LocationData struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type MsgText struct {
}

// TextEntity Type: spoiler strikethrough underline plain phone mention_name mention link italic hashtag email custom_emoji code bot_command bold bank_card
// type = text_link => + field href
// type = pre => + field language
type TextEntity struct {
	Type     string `json:"type,omitempty"`
	Text     string `json:"text,omitempty"`
	Href     string `json:"href,omitempty"`     // for type text_link
	Language string `json:"language,omitempty"` // for type pre
	UserIid  int64  `json:"user_id,omitempty"`  // for type mention_name
}

// mime_type:
//application/vnd.openxmlformats-officedocument.wordprocessingml.document
//application/x-bittorrent
//audio/flac
//audio/m4a
//audio/mp3
//audio/mpeg
//audio/mpeg3
//audio/ogg
//image/jpeg
//image/png
//text/plain
//video/mp4
//video/quicktime
