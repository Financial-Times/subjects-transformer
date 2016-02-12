package main

import (
	"github.com/pborman/uuid"
)

type subjectTransformer struct {
}

func (tr *subjectTransformer) transform(t term) subject {
	return subject{
		UUID:          uuid.NewMD5(uuid.UUID{}, []byte(t.ID)).String(),
		CanonicalName: t.CanonicalName,
		TmeIdentifier: t.ID,
		Type:          "Subject",
	}
}
