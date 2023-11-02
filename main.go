package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/ddl-hust/pied-piper/common/kafka"
	"github.com/ddl-hust/pied-piper/common/log"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// {
// 	"block_number": 18422621,
// 	"block_timestamp": "2023-10-24 20:52:59",
// 	"transaction_index": 100,
// 	"transaction_hash": "0x61eeb30f560ef4c269d3915314a4dc588aee4662a7e4ee7b16260647ae921bd2",
// 	"log_index": 234,
// 	"contract_address": "0x14fee680690900ba0cccfc76ad70fd1b95d10e16",
// 	"_from": "0x7851fa1bcadb0e7894b83c33de17bf9bb14fea5b",
// 	"_to": "0x163ad6ac78ffe40e194310faeada8f6615942d7b",
// 	"_value": "13000000000000",
// 	"pk": 3198
// }

type Transfer struct {
	From        string `json:"_from"`
	To          string `json:"_to"`
	Value       string `json:"_value"`
	Hash        string `json:"transaction_hash"`
	BlockNumber int    `json:"block_number"`
	BlockTime   string `json:"block_timestamp"`
}

var (
	PEPE_PAIR_ADDRESS = "0xa43fe16908251ee70ef74718545e4fe6c5ccec9f"
	PEPE_ADDRESS      = "0x6982508145454Ce325dDbE47a25d4ec3d2311933"
	AMOUNT_THRESHOLD  = 10000
)

func main() {

	// cfg := config.GetConf()
	// client, err := ethclient.Dial(fmt.Sprintf("https://ethereum-mainnet.s.chainbase.online/v1/%s", cfg.APIKey))
	client, err := ethclient.Dial(fmt.Sprintf("http://localhost:8545"))
	if err != nil {
		log.Error("failed to dial", zap.Error(err))
	}
	bn, _ := client.BlockNumber(context.Background())

	log.Info("[RPC] connected success", zap.Any("latest block number:", bn))

	kc, err := kafka.NewClient("ethereum.erc20.transfer")

	for {
		fetchs := kc.PollFetches(context.Background())
		records := fetchs.Records()
		for _, record := range records {
			var transfer Transfer
			json.Unmarshal(record.Value, &transfer)

			if err != nil {
				log.Error("failed to unmarshal", zap.Error(err))

			}
			if num, _ := strconv.Atoi(transfer.Value); num < AMOUNT_THRESHOLD ||
				transfer.To != PEPE_PAIR_ADDRESS {
				continue
			}
			log.Info(fmt.Sprintf("timestamp:%s,block number:%d,hash: %s, from: %s, to: %s, value: %s", transfer.BlockTime, transfer.BlockNumber, transfer.Hash, transfer.From, transfer.To, transfer.Value))

			// send tx
			cmd := exec.Command("zsh", "-c", "cast rpc anvil_impersonateAccount $USER")
			_ = cmd.Run()

			cmd = exec.Command("zsh", "-c", "cast send $ROUTER --unlocked --from $USER 0x3593564c000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000653c7bf30000000000000000000000000000000000000000000000000000000000000002090c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000547a9c189cb7587083746db00000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000006982508145454ce325ddbe47a25d4ec3d2311933000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000de0b6b3a7640000")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Error("failed to send tx", zap.Error(err))
			}
			os.Exit(1)
		}

	}
}

func filter(transfer *Transfer) bool {

	if num, _ := strconv.Atoi(transfer.Value); num < AMOUNT_THRESHOLD ||
		transfer.To != PEPE_PAIR_ADDRESS {
		return true
	}
	return false
}
