package user

import "time"

// UserProfile interface that will hold user struct.
type Profiler interface {
	GetCreationDate() string
}

func (p *Profile) GetCreationDate() string {
	return time.Now().Format("2 Jan 2006 - 15:04")
}

// User profile struct.
type Profile struct {
	ID      int
	Name    string
	Email   string
	Profile string
	Picture string
}
