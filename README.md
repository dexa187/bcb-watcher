# bcb-watcher
Simple program to watch an address for token transfers.

## Usage:
```
bcb-watcher -help
```

## Example:
```
$ bcb-watcher-linux-amd64 -wsURL=ws://localhost:9545 -contract=0xe82B7e822959B0E8131e0913ee72465be6709094 -output=json
{"From":"0xf8c29461d473daf561912ba76f441dfd8d2cc6bf","To":"0xbcb0ba1101000000000000000000000000000000","Tokens":1,"Data":"AQ=="}
{"From":"0xf8c29461d473daf561912ba76f441dfd8d2cc6bf","To":"0xbcb0ba1101000000000000000000000000000000","Tokens":1,"Data":"AQ=="}
```

## Build:
```
go build
```

## Cross Compile
```
docker pull karalabe/xgo-latest
go get github.com/karalabe/xgo
xgo --targets=linux/amd64,windows/*,darwin/amd64 ./
```