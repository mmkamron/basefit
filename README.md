# Library

Library API with CRUD operations.

### Installation with Docker

Make sure you have installed docker and it is up and running.
`docker-compose.yml` will set up a postgres database and a server.

```bash
docker-compose up
```

### Usage

A web server will start and listen on `localhost:8080`.
- To create a book, make a `POST` request to `localhost:8080/book` endpoint
  with json data that contains fields `name` and `author`.
- To retrieve all books, make a `GET` request to `localhost:8080/book` endpoint.
- To retrieve one book, make a `GET` request to `localhost:8080/book/id_number` endpoint.
- To delete a book, make a `DELETE` request to `localhost:8080/book/id_number` endpoint.
- To update a book, make a `PUT` request to `localhost:8080/book` endpoint with
  json data that contains fields `id`, `name` and `author`.

### Note

This project is in development. To contribute, please reach out to me, if you do frontend.
