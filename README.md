# Library

Library API with CRUD operations.

### Installation with Docker

Make sure you have docker and it is up and running. Currently, Dockerfile sets up a postgres database with the default username `postgres`, password and database name `library`. Then, you have to manually run a `main.go` file. In future, I will implement this in docker as well. To build a docker image, simply write:
```bash
docker build -t library .
docker run -d -p 5432:5432 library
```

### Usage

A web server will start and listen on localhost port 8080.
- To create a book, make a `POST` request to `localhost:8080/book` endpoint
  with json data that contains fields `name` and `author`. To retrieve
- To retrieve all books, make a `GET` request to `localhost:8080/book` endpoint
- To retrieve one book, make a `GET` request to `localhost:8080/book/id_number` endpoint
- To delete a book, make a `DELETE` request to `localhost:8080/book/id_number` endpoint
- To update a book, make a `PUT` request to `localhost:8080/book` endpoint with
  json data that contains fields `id`, `name` and `author`.

### Note
This project is in development. To contribute, please reach out to me, if you do frontend.
