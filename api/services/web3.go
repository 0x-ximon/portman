package services

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

const AccountAddress = "0x92572654c31cba491d7d2bc21b19349eba3ad1eb"
const ManagerAddress = "0x53993ddd5e92b9fe855ac7b48503c31c96985e31"
const DispenserAddress = "0xbb24509adf4f8aaa9582ad089364aebe281480d7"

type Web3Service struct {
	client *ethclient.Client
}

func NewWeb3Service() (*Web3Service, error) {
	rpcUrl := os.Getenv("ETHEREUM_RPC")
	if rpcUrl == "" {
		return nil, fmt.Errorf("ETHEREUM_RPC is not set")
	}

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	return &Web3Service{client: client}, nil
}
