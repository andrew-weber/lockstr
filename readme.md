<img src="lockstr.png" height="300" width="300">

# Lockstr

A password manager that uses [nostr](https://github.com/nostr-protocol/nostr).

## Usage

### Add a password

```
go run main.go add <name>
# Enter password
```

### Get a password

```
go run main.go get <name>
# Added to clipboard
```

### List all passwords

```
go run main.go list
```

## Setup

### Create Config

Add to `~/.lockstr`

```
KEY=<Private Nostr Key>
# Pick your favorite relays
RELAYS=wss://nostr.zebedee.cloud, wss://nostr-pub.wellorder.net, wss://relay.damus.io, wss://nostr.onsats.org
```
