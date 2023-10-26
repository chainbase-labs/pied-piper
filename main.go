package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ddl-hust/pied-piper/common/config"
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
)

func main() {

	cfg := config.GetConf()
	client, err := ethclient.Dial(fmt.Sprintf("https://ethereum-mainnet.s.chainbase.online/v1/%s", cfg.APIKey))
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
			if num, _ := strconv.Atoi(transfer.Value); num > 10000 &&
				transfer.To == PEPE_PAIR_ADDRESS {
				log.Info(fmt.Sprintf("timestamp:%s,block number:%d,hash: %s, from: %s, to: %s, value: %s", transfer.BlockTime, transfer.BlockNumber, transfer.Hash, transfer.From, transfer.To, transfer.Value))

			}
			// log.Info(fmt.Sprintf("hash: %s, from: %s, to: %s, value: %s", transfer.Hash, transfer.From, transfer.To, transfer.Value))

		}

	}
}
