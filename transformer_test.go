package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		term    term
		subject subject
	}{
		{"Transform term to subject", term{
			CanonicalName: "SubjectZ",
			RawID:         "Nstein_GL_AFTM_GL_164835"},
			subject{
				UUID:      "3a845a8d-944d-364e-8670-81f26434546e",
				PrefLabel: "SubjectZ",
				AlternativeIdentifiers: alternativeIdentifiers{
					TME:   []string{"TnN0ZWluX0dMX0FGVE1fR0xfMTY0ODM1-U3ViamVjdHM="},
					Uuids: []string{"3a845a8d-944d-364e-8670-81f26434546e"},
				},
				PrimaryType:   primaryType,
				TypeHierarchy: subjectTypes,
			}},
	}

	for _, test := range tests {
		expectedSubject := transformSubject(test.term, "Subjects")
		assert.Equal(test.subject, expectedSubject, fmt.Sprintf("%s: Expected subject incorrect", test.name))
	}

}
