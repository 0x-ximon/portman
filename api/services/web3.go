package services

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/0x-ximon/portman/api/contracts"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const DispenserAddress = "0x41f945b354b7a05ce10d7a14ee9a79055fd55e5f"

type Web3Service struct {
	auth *bind.TransactOpts
	conn *ethclient.Client
}

func NewWeb3Service() (*Web3Service, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Failed to get user home directory: %v", err)
	}

	sender := os.Getenv("SENDER")
	if sender == "" {
		return nil, fmt.Errorf("SENDER is not set")
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("PASSWORD is not set")
	}

	rpcUrl := os.Getenv("ETHEREUM_RPC")
	if rpcUrl == "" {
		return nil, fmt.Errorf("ETHEREUM_RPC is not set")
	}

	conn, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Ethereum RPC: %v", err)
	}

	chainID, err := conn.ChainID(context.Background())
	if err != nil {
		panic(fmt.Errorf("Failed to retrieve chain ID: %v", err))
	}

	ks := keystore.NewKeyStore(filepath.Join(home, ".ethereum", "keystores"), keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.Find(accounts.Account{Address: common.HexToAddress(sender)})
	if err != nil {
		return nil, fmt.Errorf("Failed to find account: %v", err)
	}

	bytes, err := os.ReadFile(acc.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read keystore file: %v", err)
	}

	key, err := keystore.DecryptKey(bytes, password)
	if err != nil {
		return nil, fmt.Errorf("Failed to decrypt key: %v", err)
	}

	auth := bind.NewKeyedTransactor(key.PrivateKey, chainID)
	return &Web3Service{conn: conn, auth: auth}, nil
}

func (w *Web3Service) SendTransaction() (string, error) {
	return "", nil
}

func (w *Web3Service) FundAccount(ctx context.Context, receiver accounts.Account) (*common.Hash, error) {
	dispenser := contracts.NewDispenser()
	instance := dispenser.Instance(w.conn, common.HexToAddress(DispenserAddress))

	tx, err := bind.Transact(instance, w.auth, dispenser.PackDispense("USD", receiver.Address, big.NewInt(10000e6)))
	if err != nil {
		return nil, fmt.Errorf("Failed to transact: %v", err)
	}

	r, err := bind.WaitMined(ctx, w.conn, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("Failed to wait for transaction to be mined: %v", err)
	}

	return &r.TxHash, nil
}
