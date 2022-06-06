package session_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/instill-ai/usage-client/session"
)

func TestNormSessionEdition_Valid(t *testing.T) {
	editions := []string{"local-ce", "local-ce:dev"}
	for _, e := range editions {
		normEdition := session.NormSessionEdition(e)

		require.Equal(t, e, normEdition)
	}
}

func TestNormSessionEdition_Invalid(t *testing.T) {
	edition := "invalid-edition"
	normEdition := session.NormSessionEdition(edition)

	require.Equal(t, session.DefaultSessionEdition, normEdition)
}
