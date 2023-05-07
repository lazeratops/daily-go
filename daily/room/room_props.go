package room

import "time"

// Props represents common properties the user can set
// when creating or updating a room.
// This does not represent _all_ properties supported
// by Daily rooms.
type Props struct {
	// Exp should be a Unix timestamp, but we'll provide
	// some helper methods to let caller work with time.Time
	// as well
	Exp             int64 `json:"exp,omitempty"`
	MaxParticipants int   `json:"max_participants,omitempty"`
	StartAudioOff   bool  `json:"start_audio_off"`
	StartVideoOff   bool  `json:"start_video_off"`
}

func GetRoomPropsKeys() []string {
	return []string{"exp", "max_participants", "start_audio_off", "start_video_off"}
}

// SetExpiry sets the room expiry as a Unix timestamp
func (p *Props) SetExpiry(expiry time.Time) {
	p.Exp = expiry.Unix()
}

// GetExpiry retrieves the room expiry
func (p *Props) GetExpiry() time.Time {
	return time.Unix(p.Exp, 0)
}
