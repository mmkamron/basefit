# Library

Library API with CRUD operations.

Note: For testing to be easier, credentials to the database are given as constants.

### Installation with Docker

Make sure you have installed docker and it is up and running.
`docker-compose.yml` will set up a postgres database and a server. Simply run:

```bash
docker-compose up
```

### Usage

A web server will start and listen on `localhost:8080` listing all the books in the library. You have to login to view the table of books. By default, there are 3 books in the database.
- To retrieve one book, go to the `localhost:8080/book/id_number`.
- To delete a book, make a `DELETE` request to `localhost:8080/book/id_number` endpoint.
- To update a book, make a `PUT` request to `localhost:8080/book` endpoint with
  json data that contains fields `id`, `name` and `author`.

This project is under development.
