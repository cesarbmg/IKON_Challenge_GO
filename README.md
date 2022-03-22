# IKON_Challenge_GO

Project Ikon_Challenge with GO

## Imports package

The app need import packages with "go get"

- go install google.golang.org/grpc@latest
- go install github.com/cesarbmg/IKON_Challenge_GO/grpc@latest

## External files

The app need files of input

- challenge.in

Example Sintaxis:

- 20                          // First Line is Capacity of Device
- (1, 7), (2, 14), (3, 8)     // Second Line is Tasks of Foreground: (IdTask, Resources)
- (1, 14), (2, 5), (3, 10)    // Third Line is Tasks of Background: (IdTask, Resources)

## Command run

The app execute 3 files "go run" in same time

- go run REST\Server_REST.go
- go run gRPC\Server_gRPC.go
- go run Client.go
