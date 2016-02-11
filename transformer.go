package main

type subjectTransformer struct {
}

func (t *subjectTransformer) transform(term Term) Subject {
	return Subject{
		UUID:          NewNameUUIDFromBytes([]byte(term.Id)).String(),
		CanonicalName: term.CanonicalName,
		TmeIdentifier: term.Id,
		Type:          "Subject",
	}
}
