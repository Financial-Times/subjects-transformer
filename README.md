# subjects-transformer

[![Circle CI](https://circleci.com/gh/Financial-Times/subjects-transformer/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/subjects-transformer/tree/master)

Retrieves Subjects taxonomy from TME vie the structure service and transforms the subjects to the internal UP json model.
The service exposes endpoints for getting all the subjects and for getting subject by uuid.

# Usage
`go get github.com/Financial-Times/subjects-transformer`

`$GOPATH/bin/subjects-transformer --port=8080 -base-url="http://localhost:8080/transformers/subjects/" -structure-service-base-url="http://metadata.internal.ft.com:83" -structure-service-username="user" -structure-service-password="pass" -structure-service-principal-header="app-preditor"`
```
export|set PORT=8080
export|set BASE_URL="http://localhost:8080/transformers/subjects/"
export|set STRUCTURE_SERVICE_BASE_URL="http://metadata.internal.ft.com:83"
export|set STRUCTURE_SERVICE_USERNAME="user"
export|set STRUCTURE_SERVICE_PASSWORD="pass"
export|set PRINCIPAL_HEADER="app-preditor"
$GOPATH/bin/subjects-transformer
```

With Docker:

`docker build -t coco/subjects-transformer .`

`docker run -ti --env BASE_URL=<base url> --env STRUCTURE_SERVICE_BASE_URL=<structure service url> --env STRUCTURE_SERVICE_USERNAME=<user> --env STRUCTURE_SERVICE_PASSWORD=<pass> --env PRINCIPAL_HEADER=<header> coco/subjects-transformer`
