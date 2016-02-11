package main

type subjectTransformer struct {
}

func (tr *subjectTransformer) transform(t term) subject {
	return subject{
		UUID:          NewNameUUIDFromBytes([]byte(t.ID)).String(),
		CanonicalName: t.CanonicalName,
		TmeIdentifier: t.ID,
		Type:          "Subject",
	}
}
