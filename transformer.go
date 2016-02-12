package main

import (
	"github.com/pborman/uuid"
)

func transformSubject(t term) subject {
	return subject{
		UUID:          uuid.NewMD5(uuid.UUID{}, []byte(t.ID)).String(),
		CanonicalName: t.CanonicalName,
		TmeIdentifier: t.ID,
		Type:          "Subject",
	}
}
