package service

import (
	"github.com/Financial-Times/subjects-transformer/model"
)


type SubjectTransformer struct {
}

func (t *SubjectTransformer) transform(term model.Term) model.Subject {
	return model.Subject{
		UUID: NewNameUUIDFromBytes([]byte(term.Id)).String(),
		CanonicalName: term.CanonicalName,
		TmeIdentifier: term.Id,
		Type: "Subject",
	}
}
