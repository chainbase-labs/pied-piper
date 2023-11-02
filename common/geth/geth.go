package geth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ddl-hust/pied-piper/common/config"
	"github.com/ddl-hust/pied-piper/common/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"go.uber.org/zap"
)

type GethClient struct {
	Client *ethclient.Client
}

func NewGethClient(cfg *config.Config) (*GethClient, error) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Error("failed to dial", zap.Error(err))
		return nil, err
	}
	log.Info("geth client connected", zap.String("url", "http://localhost:8545"))
	return &GethClient{
		Client: client,
	}, nil
}

// 1.load private key
// 2.make up calldata
// 3.sign
// 4.send
func (geth *GethClient) SendTx(cfg *config.Config, to string) error {
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("failed to assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := geth.Client.PendingNonceAt(context.Background(), fromAddress)

	value := big.NewInt(1)     // in wei (1 eth)
	gasLimit := uint64(210000) // in units
	gasPrice, err := geth.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error("failed to suggest gas price", zap.Error(err))
		return err
	}

	toAddress := common.HexToAddress(to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := geth.Client.NetworkID(context.Background())
	if err != nil {
		log.Error("failed to get network id", zap.Error(err))
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Error("failed to sign tx", zap.Error(err))
		return err
	}

	err = geth.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Error("failed to send tx", zap.Error(err))
		return err
	}
	bn, err := geth.Client.BlockNumber(context.Background())
	log.Info("block number:", zap.Any("latest block number:", bn))

	log.Info(fmt.Sprintf("tx sent: %s", signedTx.Hash().Hex()))
	return nil
}
