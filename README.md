Reddys is a simple Redis server written in Go.

So far I've implemented:
1. Commands: `GET/SET/MSET`, `HGET/HSET`, `INCR/DECR/INCRBY/DECRBY`.
2. Persistance of the DB via AOF file.

Run via `go run github.com/dsocolobsky/reddys/cmd/server`.

Then you may connect via `redis-cli`
