package room

import "time"

// RoomProps represents common properties the user can set
// when creating or updating a room.
// This does not represent _all_ properties supported
// by Daily rooms.
type RoomProps struct {
	// Exp should be a Unix timestamp, but we'll provide
	// some helper methods to let caller work with time.Time
	// as well
	Exp             int64 `json:"exp,omitempty"`
	MaxParticipants int   `json:"max_participants,omitempty"`
}

func GetRoomPropsKeys() []string {
	return []string{"exp", "max_participants"}
}

// SetExpiry sets the room expiry as a Unix timestamp
func (p *RoomProps) SetExpiry(expiry time.Time) {
	p.Exp = expiry.Unix()
}

// GetExpiry retrieves the room expiry
func (p *RoomProps) GetExpiry() time.Time {
	return time.Unix(p.Exp, 0)
}
