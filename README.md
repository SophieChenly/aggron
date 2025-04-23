

# AGGRON
![Untitled_Artwork - 2025-04-14T235541 516](https://github.com/user-attachments/assets/a3288b35-7b61-4647-8247-a4182205c864)

## Description
- This is a discord bot that encrypts files sent by users and also bypasses file size limits. 
- Built with Go, AWS, Discord API, and Passage by 1Password (auth).

## Setup and Run
Install Go: https://go.dev/doc/install

1. `go mod download`
2. `go run main.go` (fast run)

OR if you want to build and run the binary (production)

3. `go build`
4. `./aggron.exe` (or whatever executable)


## Project Structure

├── internal/
│   ├── api/
│   │   └── handlers.go      # where the endpoints controller goes
│   ├── bot/
│   │   └── discord.go       # where the discord bot handler goes
│   ├── models/
│   │   └── types.go
│   └── services/            # handles business logic (i.e. 1password, file uploads, encryption classes)
│       ├── service.go
│