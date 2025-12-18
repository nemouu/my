package ai

// This file will contain session management logic for maintaining
// conversation context, loaded files, and token tracking on a
// per-directory basis.
//
// TODO: Implement session management as described in TODO.md Milestone 3 "Session Management"
// - Define Session struct (Directory, Messages, LoadedFiles, TotalTokens, etc.)
// - Generate session file path from directory hash (SHA256)
// - Implement LoadSession(dir string) (*Session, error)
// - Implement SaveSession(session *Session) error
// - Implement approximate token counting

// Session represents a conversation session for a specific directory
type Session struct {
	// TODO: Add fields for Directory, Messages, LoadedFiles, TotalTokens, LastActive, etc.
}

// LoadSession loads or creates a session for the given directory
func LoadSession(dir string) (*Session, error) {
	// TODO: Load session from ~/.config/my/sessions/{hash}.json
	return nil, nil
}

// SaveSession persists the session to disk
func (s *Session) SaveSession() error {
	// TODO: Save session to JSON file
	return nil
}

// AddMessage adds a message to the session history
func (s *Session) AddMessage(role, content string) {
	// TODO: Append message and update token count
}

// LoadFile loads a file's contents into the session context
func (s *Session) LoadFile(path string) error {
	// TODO: Read file and add to LoadedFiles map
	return nil
}
