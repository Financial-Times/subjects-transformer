# subjects-transformer

[![Circle CI](https://circleci.com/gh/Financial-Times/subjects-transformer/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/subjects-transformer/tree/master)

Retrieves Subjects taxonomy from TME vie the structure service and transforms the subjects to the internal UP json model.
The service exposes endpoints for getting all the subjects and for getting subject by uuid.

# Usage
`go get github.com/Financial-Times/subjects-transformer`

`subjects-transformer -baseUrl=http://localhost:8080/transformers/subjects/ -structureServiceBaseUrl=http://metadata.internal.ft.com:83 -structureServiceUsername=user -structureServicePassword=pass -structureServicePrincipalHeader=app-preditor`

With Docker:

`docker build -t coco/subjects-transformer .`

`docker run -ti --env BASE_URL=<base url> --env STRUCTURE_SERVICE_BASE_URL=<structure service url> --env USER=<user> --env PASS=<pass> --env PRINCIPAL_HEADER=<header> coco/subjects-transformer`
