# Sock
Sock's minimal lightweight core server written in Go.

## Flags
- `limit` (int): Maximum read limit (default 1kb).
- `ping` (int): Wait for ping in x second.
- `port` (string): Set sock sevre port (default `process.env.PORT` | 3000).
- `token` (string): Set sock token to only access who knows token (default `demo`).

## Example(s)
- `./sock`
- `./sock -port=8080`
- `./sock -limit=10485760 -token=secret-key`