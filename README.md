# lexportd

Snapshot and export LXC containers using LXD socket API (**L**e**X**port**D**).

## Usage

**Note:** The LXD socket may be accessible to `root` only so it is preferable to run `go build` instead of `go run`.

```bash
./lexportd [args...]
```

### Arguments

| Argument |             Description              |          Default          |
| :------: | :----------------------------------: | :-----------------------: |
| `-sock`  |          Path to LXD socket          |        `REQUIRED`         |
|  `-out`  | Location to write exported snapshots | Current working directory |
