package session

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	usagePB "github.com/instill-ai/protogen-go/vdp/usage/v1alpha"
)

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

var (
	DefaultSessionEdition = "local-ce:dev"
	ValidSessionEditions  = []string{"local-ce", DefaultSessionEdition}
)

// NormSessionEdition normalises a session edition input.
// If it's not valid, then `DefaultSessionEdition` is returned, or else the original edition input is returned.
func NormSessionEdition(edition string) string {

	if valid := contains(ValidSessionEditions, edition); valid {
		return edition
	} else {
		return DefaultSessionEdition
	}
}

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
