.PHONY: migrate
migrate:
	migrate \
		-source file://db/changelog \
		-database "postgres://opensaga@localhost:5432/opensaga?sslmode=disable&search_path=maintenance,opensaga,public&x-migrations-table=changelog" \
		up
