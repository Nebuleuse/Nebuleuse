#Nebuleuse
Nebuleuse is a web API for game developpers to use to integrate Stats, Achievements, Matchmaking, Inventory and more to their games. This repository is focused on the server Backend written in Go.
#Clients
Currently only one client exists for [C++] but it's easy to port and create a client in another language.
#Installation and Building
- Install [Go]
- go get github.com/go-sql-diver/mysql
- go get github.com/nu7hatch/gouuid
- go get github.com/gorilla/mux
- go install

#API
Nebuleuse uses a REST-like API for its communications, the full API documentation is avialable [here][ApiWiki].

[C++]:https://github.com/Orygin/NebuleuseCppClient
[Go]:https://golang.org/doc/install
[ApiWiki]:https://github.com/Orygin/Nebuleuse/wiki/API