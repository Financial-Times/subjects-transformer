package main

type SubjectTransformer struct {
}

func (t *SubjectTransformer) transform(term Term) Subject {
	return Subject{
		UUID:          NewNameUUIDFromBytes([]byte(term.Id)).String(),
		CanonicalName: term.CanonicalName,
		TmeIdentifier: term.Id,
		Type:          "Subject",
	}
}
