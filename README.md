# gofamily

This is a small hack to make a web-service that displays
genealogy data imported from Scion to a postgres database
in s small go-lang application.

## Design overview

Using docker-compose to run locally: two services, postgres
the database and service the application.

The applications have a main class gofamily that displays (well
"hello world") and also there is a simple database importer that
reads likely a subset of Scion.xml, hardcoded to Arnholm.sgx.  