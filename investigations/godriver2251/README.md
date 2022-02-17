`go get` with cloud-1.7.1 branch is considered a downgrade.

```
(.venv) kevin.albertson@M-C02D44JBML85 godriver2251 % go get go.mongodb.org/mongo-driver@v1.7.2
go: downloading go.mongodb.org/mongo-driver v1.7.2
go get: upgraded go.mongodb.org/mongo-driver v1.7.1 => v1.7.2
(.venv) kevin.albertson@M-C02D44JBML85 godriver2251 % go get go.mongodb.org/mongo-driver@cloud-1.7.1
go: downloading go.mongodb.org/mongo-driver v1.4.0-beta2.0.20220104204243-ea4595952704
go get: downgraded go.mongodb.org/mongo-driver v1.7.2 => v1.4.0-beta2.0.20220104204243-ea4595952704
(.venv) kevin.albertson@M-C02D44JBML85 godriver2251 % go mod tidy
go: downloading golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
```