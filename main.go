package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"strings"

	token "github.com/dexa187/bcb-watcher/contracts"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
	Data   []byte
}

func main() {

	contractPtr := flag.String("contract", "0x392b695c3da2D86ec5284fAe2D152231876a3548", "Address of BCB Contract")
	walletPtr := flag.String("wallet", "0xbCb0Ba1101000000000000000000000000000000", "Address of the Wallet to Watch")
	wsURLPtr := flag.String("wsURL", "wss://dai-trace-ws.blockscout.com/ws", "Websocket URL of blockchain node")
	outputPtr := flag.String("output", "kv", "Output mode kv or json")

	flag.Parse()

	client, err := ethclient.Dial(*wsURLPtr)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(*contractPtr)
	toAddressTopic := common.HexToAddress(*walletPtr).Hash()
	eventTopic := crypto.Keccak256Hash([]byte("TransferWithData(address,address,uint256,bytes)"))

	// This contract abi is not the full BCB contract.  Just the events we are looking for
	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	// The Addresses contains limits the subscription to just the bcb contract address
	// Topics are hashes of the events arguments.
	// In this case the first topic is the hash of the event name and arg types
	// The second topic is the FROM address which is nil meaning any
	// The third topic is the TO address which is the address we want to watch
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{eventTopic}, nil, {toAddressTopic}},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:

			// We are only watching for TransferWithData events so all logs will be the same event type.
			var transferEvent LogTransfer

			err := contractAbi.Unpack(&transferEvent, "TransferWithData", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			switch *outputPtr {
			case "kv":
				fmt.Printf("From=%s To=%s Tokens=%s\n Data=%s", transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Tokens.String(), string(transferEvent.Data))
			case "json":
				out, _ := json.Marshal(transferEvent)
				fmt.Printf(string(out))
			}

		}
	}
}
