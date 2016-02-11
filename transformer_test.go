package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	transformer := SubjectTransformer{}
	tests := []struct {
		name    string
		term    Term
		subject Subject
	}{
		{"Trasform term to subject", Term{CanonicalName: "Metals Markets", Id: "MTE3-U3ViamVjdHM="}, Subject{UUID: "bba39990-c78d-3629-ae83-808c333c6dbc", CanonicalName: "Metals Markets", TmeIdentifier: "MTE3-U3ViamVjdHM=", Type: "Subject"}},
	}

	for _, test := range tests {
		expectedSubject := transformer.transform(test.term)
		assert.Equal(test.subject, expectedSubject, fmt.Sprintf("%s: Expected subject incorrect", test.name))
	}

}
