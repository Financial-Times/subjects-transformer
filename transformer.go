package main

type subjectTransformer struct {
}

func (tr *subjectTransformer) transform(t term) subject {
	return subject{
		UUID:          NewNameUUIDFromBytes([]byte(t.Id)).String(),
		CanonicalName: t.CanonicalName,
		TmeIdentifier: t.Id,
		Type:          "Subject",
	}
}
