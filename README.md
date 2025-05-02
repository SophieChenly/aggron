

# AGGRON
![Untitled_Artwork - 2025-04-14T235541 516](https://github.com/user-attachments/assets/a3288b35-7b61-4647-8247-a4182205c864)

## Description
- This is a discord bot that encrypts files sent by users and also bypasses file size limits. 
- Built with Go, AWS, Discord API, and Passage by 1Password (auth).

## Setup and Run
Install Go: https://go.dev/doc/install

1. `go mod download`
2. `docker pull redis`
3. `docker run -p 6379:6379 -d redis`
4. `go run main.go` (fast run)

OR if you want to build and run the binary (production)

4. `go build`
5. `./aggron.exe` (or whatever executable)
