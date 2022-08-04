# Library

Library API with CRUD operations.

### Installation with Docker

Make sure you have installed docker and it is up and running.
`docker-compose.yml` will set up a postgres database and a server. Simply run:

```bash
docker-compose up
```

### Usage

A web server will start and listen on `localhost:8080` listing all the books in the library. By default, there are 3 books in the database.
- To create a book, make a `POST` request to `localhost:8080/book` endpoint
  with json data that contains fields `name` and `author`.
- To retrieve one book, make a `GET` request to `localhost:8080/book/id_number` endpoint.
- To delete a book, make a `DELETE` request to `localhost:8080/book/id_number` endpoint.
- To update a book, make a `PUT` request to `localhost:8080/book` endpoint with
  json data that contains fields `id`, `name` and `author`.

This project is in development.
