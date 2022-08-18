# Async Chess Server

This is a server for an asynchronous game of Chess. The client is written in Rust, but I wanted to write a server in Go to get more Go experience. If you're running this yourself, make sure to add a `.env` file with a variable for `DB_PASSWORD` and to change the server ip in `main.go`


The API currently relies on honesty in terms of who needs to take a turn, and where pieces move to (somewhat like real chess...). Currently no support for promotion

## Endpoints

The server exposes a REST API using the gin go library with the following endpoints:

### GET:

 - `/games/:id` - pass in an integer id to get all of the pieces involved in that game. Will 500 if no pieces. Will return 200 with pieces, or a 208 if the pieces haven't changed.

 ### POST:

 - `/newgame` - pass in an integer id to create a new game, and this will overwrite the old game
 - `/deletegame` - pass in an integer id to delete a game
 - `/movepiece` - pass in a json struct of all integers: `id, x, y, nx, ny` where all are > 0, and `x,y,nx,ny` are < 8. All fields must be set for validation purposes, as much as Go loves the default values.

## TODO:
 - GUI desktop client
 - Web client
 - validation to ensure the same users are playing the same games, maybe via cookies/passwords
 - TUI for running the game, whilst somehow preserving the `gin` logs
 - validation for movements of chess pieces, eg. knights jump, pawns take diagonally
 - validation for moves made in correct order
 - validation for when the game is done
 - support for weird rules like castling and promotion
