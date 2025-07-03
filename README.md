# corvid

a dbus notification server that outputs json. intended for use with [eww](https://github.com/elkowar/eww).

## subcommands
- `server`: run the server
- `dismiss <id>`: dismiss specific notification
- `dismiss-all`: dismiss all notifications

## environment variables
- `CORVID_DEFAULT_EXPIRATION`: default notification expiration in ms. `-1` = never (default: `5000`)
- `CORVID_SORT_DIRECTION`: notification sort direction, either `NEWEST_FIRST` or `OLDEST_FIRST` (default: `NEWEST_FIRST`)

## building

### building with nix (recommended)
1. clone
    ```bash
    git clone https://github.com/CartConnoisseur/corvid
    cd corvid
    ```
2. build
    ```bash
    nix build
    ```

### building with go
#### prerequisites:
- go

1. clone
    ```bash
    git clone https://github.com/CartConnoisseur/corvid
    cd corvid
    ```
2. build
    ```bash
    go mod tidy
    go build -o corvid ./...
    ```
