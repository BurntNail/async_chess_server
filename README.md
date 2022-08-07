# Async Chess Server

This is a server for an asynchronous game of Chess. The client is written in Rust, but I wanted to write a server in Go to get more Go experience. If you're running this yourself, make sure to add a `.env` file with a variable for `DB_PASSWORD` and to change the server ip in `main.go`


The API currently relies on honesty in terms of who needs to take a turn, and where pieces move to (somewhat like real chess...). Currently no support for promotion

## Endpoints

The server exposes a REST API using the gin go library with the following endpoints:

### GET:

 - `/games/:id` - pass in an integer id to get all of the pieces involved in that game. Will 500 if no pieces.

 ### POST:

 - `/newgame` - pass in an integer id to create a new game, and this will overwrite the old game
 - `/deletegame` - pass in an integer id to delete a game
 - `/movepiece`