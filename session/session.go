package session

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	usagePB "github.com/instill-ai/protogen-go/core/usage/v1beta"
)

// Session is a new type of usagePB.session
type Session usagePB.Session

// Hash converts the usage data into base64 encoded checksum
func (s *Session) Hash() (string, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(s)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(buf.Bytes())
	return base64.URLEncoding.EncodeToString(sum[:]), nil
}

// Expired checks whether the usage data is expired
func (s *Session) Expired() bool {
	return time.Since(s.ReportTime.AsTime()) > time.Minute
}
