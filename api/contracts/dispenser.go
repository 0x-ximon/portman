// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// DispenserMetaData contains all meta data concerning the Dispenser contract.
var DispenserMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"assets\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractPortmanAsset\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createAsset\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimal\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dispense\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AssetAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AssetNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "Dispenser",
	Bin: "0x608060405234801561000f575f5ffd5b50604051612647380380612647833981810160405281019061003191906101d7565b805f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100999190610211565b60405180910390fd5b6100b1816100b860201b60201c565b505061022a565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101a68261017d565b9050919050565b6101b68161019c565b81146101c0575f5ffd5b50565b5f815190506101d1816101ad565b92915050565b5f602082840312156101ec576101eb610179565b5b5f6101f9848285016101c3565b91505092915050565b61020b8161019c565b82525050565b5f6020820190506102245f830184610202565b92915050565b612410806102375f395ff3fe608060405234801561000f575f5ffd5b5060043610610060575f3560e01c806330568c3414610064578063715018a614610080578063859362281461008a5780638da5cb5b146100ba578063da7b7ce3146100d8578063f2fde38b14610108575b5f5ffd5b61007e600480360381019061007991906107c8565b610124565b005b610088610240565b005b6100a4600480360381019061009f9190610834565b610253565b6040516100b191906108d6565b60405180910390f35b6100c261029b565b6040516100cf91906108fe565b60405180910390f35b6100f260048036038101906100ed919061094d565b6102c2565b6040516100ff91906108fe565b60405180910390f35b610122600480360381019061011d91906109d5565b61040e565b005b61012c610492565b5f60018460405161013d9190610a52565b90815260200160405180910390205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036101d2576040517f470cbf4700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff16637b1837de84846040518363ffffffff1660e01b815260040161020d929190610a77565b5f604051808303815f87803b158015610224575f5ffd5b505af1158015610236573d5f5f3e3d5ffd5b5050505050505050565b610248610492565b6102515f610519565b565b6001818051602081018201805184825260208301602085012081835280955050505050505f915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f6102cb610492565b5f73ffffffffffffffffffffffffffffffffffffffff166001856040516102f29190610a52565b90815260200160405180910390205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461036d576040517fdc0d0aab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8484843060405161037e906105e1565b61038b9493929190610af5565b604051809103905ff0801580156103a4573d5f5f3e3d5ffd5b509050806001866040516103b89190610a52565b90815260200160405180910390205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550809150509392505050565b610416610492565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610486575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161047d91906108fe565b60405180910390fd5b61048f81610519565b50565b61049a6105da565b73ffffffffffffffffffffffffffffffffffffffff166104b861029b565b73ffffffffffffffffffffffffffffffffffffffff1614610517576104db6105da565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161050e91906108fe565b60405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f33905090565b61189480610b4783390190565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61064d82610607565b810181811067ffffffffffffffff8211171561066c5761066b610617565b5b80604052505050565b5f61067e6105ee565b905061068a8282610644565b919050565b5f67ffffffffffffffff8211156106a9576106a8610617565b5b6106b282610607565b9050602081019050919050565b828183375f83830152505050565b5f6106df6106da8461068f565b610675565b9050828152602081018484840111156106fb576106fa610603565b5b6107068482856106bf565b509392505050565b5f82601f830112610722576107216105ff565b5b81356107328482602086016106cd565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6107648261073b565b9050919050565b6107748161075a565b811461077e575f5ffd5b50565b5f8135905061078f8161076b565b92915050565b5f819050919050565b6107a781610795565b81146107b1575f5ffd5b50565b5f813590506107c28161079e565b92915050565b5f5f5f606084860312156107df576107de6105f7565b5b5f84013567ffffffffffffffff8111156107fc576107fb6105fb565b5b6108088682870161070e565b935050602061081986828701610781565b925050604061082a868287016107b4565b9150509250925092565b5f60208284031215610849576108486105f7565b5b5f82013567ffffffffffffffff811115610866576108656105fb565b5b6108728482850161070e565b91505092915050565b5f819050919050565b5f61089e6108996108948461073b565b61087b565b61073b565b9050919050565b5f6108af82610884565b9050919050565b5f6108c0826108a5565b9050919050565b6108d0816108b6565b82525050565b5f6020820190506108e95f8301846108c7565b92915050565b6108f88161075a565b82525050565b5f6020820190506109115f8301846108ef565b92915050565b5f60ff82169050919050565b61092c81610917565b8114610936575f5ffd5b50565b5f8135905061094781610923565b92915050565b5f5f5f60608486031215610964576109636105f7565b5b5f84013567ffffffffffffffff811115610981576109806105fb565b5b61098d8682870161070e565b935050602084013567ffffffffffffffff8111156109ae576109ad6105fb565b5b6109ba8682870161070e565b92505060406109cb86828701610939565b9150509250925092565b5f602082840312156109ea576109e96105f7565b5b5f6109f784828501610781565b91505092915050565b5f81519050919050565b5f81905092915050565b8281835e5f83830152505050565b5f610a2c82610a00565b610a368185610a0a565b9350610a46818560208601610a14565b80840191505092915050565b5f610a5d8284610a22565b915081905092915050565b610a7181610795565b82525050565b5f604082019050610a8a5f8301856108ef565b610a976020830184610a68565b9392505050565b5f82825260208201905092915050565b5f610ab882610a00565b610ac28185610a9e565b9350610ad2818560208601610a14565b610adb81610607565b840191505092915050565b610aef81610917565b82525050565b5f6080820190508181035f830152610b0d8187610aae565b90508181036020830152610b218186610aae565b9050610b306040830185610ae6565b610b3d60608301846108ef565b9594505050505056fe60a060405234801561000f575f5ffd5b506040516118943803806118948339818101604052810190610031919061038d565b8084848160039081610043919061064a565b508060049081610053919061064a565b5050505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100c6575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100bd9190610728565b60405180910390fd5b6100d5816100ed60201b60201c565b508160ff1660808160ff168152505050505050610741565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61020f826101c9565b810181811067ffffffffffffffff8211171561022e5761022d6101d9565b5b80604052505050565b5f6102406101b0565b905061024c8282610206565b919050565b5f67ffffffffffffffff82111561026b5761026a6101d9565b5b610274826101c9565b9050602081019050919050565b8281835e5f83830152505050565b5f6102a161029c84610251565b610237565b9050828152602081018484840111156102bd576102bc6101c5565b5b6102c8848285610281565b509392505050565b5f82601f8301126102e4576102e36101c1565b5b81516102f484826020860161028f565b91505092915050565b5f60ff82169050919050565b610312816102fd565b811461031c575f5ffd5b50565b5f8151905061032d81610309565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61035c82610333565b9050919050565b61036c81610352565b8114610376575f5ffd5b50565b5f8151905061038781610363565b92915050565b5f5f5f5f608085870312156103a5576103a46101b9565b5b5f85015167ffffffffffffffff8111156103c2576103c16101bd565b5b6103ce878288016102d0565b945050602085015167ffffffffffffffff8111156103ef576103ee6101bd565b5b6103fb878288016102d0565b935050604061040c8782880161031f565b925050606061041d87828801610379565b91505092959194509250565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061047757607f821691505b60208210810361048a57610489610433565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f600883026104ec7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826104b1565b6104f686836104b1565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61053a6105356105308461050e565b610517565b61050e565b9050919050565b5f819050919050565b61055383610520565b61056761055f82610541565b8484546104bd565b825550505050565b5f5f905090565b61057e61056f565b61058981848461054a565b505050565b5f5b828110156105af576105a45f828401610576565b600181019050610590565b505050565b601f8211156106025782821115610601576105ce81610490565b6105d7836104a2565b6105e0856104a2565b60208610156105ed575f90505b8083016105fc8284038261058e565b505050505b5b505050565b5f82821c905092915050565b5f6106225f1984600802610607565b1980831691505092915050565b5f61063a8383610613565b9150826002028217905092915050565b61065382610429565b67ffffffffffffffff81111561066c5761066b6101d9565b5b6106768254610460565b6106818282856105b4565b5f60209050601f8311600181146106b2575f84156106a0578287015190505b6106aa858261062f565b865550610711565b601f1984166106c086610490565b5f5b828110156106e7578489015182556001820191506020850194506020810190506106c2565b868310156107045784890151610700601f891682610613565b8355505b6001600288020188555050505b505050505050565b61072281610352565b82525050565b5f60208201905061073b5f830184610719565b92915050565b60805161113b6107595f395f610385015261113b5ff3fe608060405234801561000f575f5ffd5b50600436106100cd575f3560e01c8063715018a61161008a57806395d89b411161006457806395d89b41146101ff578063a9059cbb1461021d578063dd62ed3e1461024d578063f2fde38b1461027d576100cd565b8063715018a6146101bb5780637b1837de146101c55780638da5cb5b146101e1576100cd565b806306fdde03146100d1578063095ea7b3146100ef57806318160ddd1461011f57806323b872dd1461013d578063313ce5671461016d57806370a082311461018b575b5f5ffd5b6100d9610299565b6040516100e69190610db4565b60405180910390f35b61010960048036038101906101049190610e65565b610329565b6040516101169190610ebd565b60405180910390f35b61012761034b565b6040516101349190610ee5565b60405180910390f35b61015760048036038101906101529190610efe565b610354565b6040516101649190610ebd565b60405180910390f35b610175610382565b6040516101829190610f69565b60405180910390f35b6101a560048036038101906101a09190610f82565b6103a9565b6040516101b29190610ee5565b60405180910390f35b6101c36103ee565b005b6101df60048036038101906101da9190610e65565b610401565b005b6101e9610417565b6040516101f69190610fbc565b60405180910390f35b61020761043f565b6040516102149190610db4565b60405180910390f35b61023760048036038101906102329190610e65565b6104cf565b6040516102449190610ebd565b60405180910390f35b61026760048036038101906102629190610fd5565b6104f1565b6040516102749190610ee5565b60405180910390f35b61029760048036038101906102929190610f82565b610573565b005b6060600380546102a890611040565b80601f01602080910402602001604051908101604052809291908181526020018280546102d490611040565b801561031f5780601f106102f65761010080835404028352916020019161031f565b820191905f5260205f20905b81548152906001019060200180831161030257829003601f168201915b5050505050905090565b5f5f6103336105f7565b90506103408185856105fe565b600191505092915050565b5f600254905090565b5f5f61035e6105f7565b905061036b858285610610565b6103768585856106a3565b60019150509392505050565b5f7f0000000000000000000000000000000000000000000000000000000000000000905090565b5f5f5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b6103f6610793565b6103ff5f61081a565b565b610409610793565b61041382826108dd565b5050565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60606004805461044e90611040565b80601f016020809104026020016040519081016040528092919081815260200182805461047a90611040565b80156104c55780601f1061049c576101008083540402835291602001916104c5565b820191905f5260205f20905b8154815290600101906020018083116104a857829003601f168201915b5050505050905090565b5f5f6104d96105f7565b90506104e68185856106a3565b600191505092915050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b61057b610793565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036105eb575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016105e29190610fbc565b60405180910390fd5b6105f48161081a565b50565b5f33905090565b61060b838383600161095c565b505050565b5f61061b84846104f1565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81101561069d578181101561068e578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161068593929190611070565b60405180910390fd5b61069c84848484035f61095c565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610713575f6040517f96c6fd1e00000000000000000000000000000000000000000000000000000000815260040161070a9190610fbc565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610783575f6040517fec442f0500000000000000000000000000000000000000000000000000000000815260040161077a9190610fbc565b60405180910390fd5b61078e838383610b2b565b505050565b61079b6105f7565b73ffffffffffffffffffffffffffffffffffffffff166107b9610417565b73ffffffffffffffffffffffffffffffffffffffff1614610818576107dc6105f7565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161080f9190610fbc565b60405180910390fd5b565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361094d575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016109449190610fbc565b60405180910390fd5b6109585f8383610b2b565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036109cc575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016109c39190610fbc565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610a3c575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401610a339190610fbc565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015610b25578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610b1c9190610ee5565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610b7b578060025f828254610b6f91906110d2565b92505081905550610c49565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610c04578381836040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401610bfb93929190611070565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610c90578060025f8282540392505081905550610cda565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610d379190610ee5565b60405180910390a3505050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f610d8682610d44565b610d908185610d4e565b9350610da0818560208601610d5e565b610da981610d6c565b840191505092915050565b5f6020820190508181035f830152610dcc8184610d7c565b905092915050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610e0182610dd8565b9050919050565b610e1181610df7565b8114610e1b575f5ffd5b50565b5f81359050610e2c81610e08565b92915050565b5f819050919050565b610e4481610e32565b8114610e4e575f5ffd5b50565b5f81359050610e5f81610e3b565b92915050565b5f5f60408385031215610e7b57610e7a610dd4565b5b5f610e8885828601610e1e565b9250506020610e9985828601610e51565b9150509250929050565b5f8115159050919050565b610eb781610ea3565b82525050565b5f602082019050610ed05f830184610eae565b92915050565b610edf81610e32565b82525050565b5f602082019050610ef85f830184610ed6565b92915050565b5f5f5f60608486031215610f1557610f14610dd4565b5b5f610f2286828701610e1e565b9350506020610f3386828701610e1e565b9250506040610f4486828701610e51565b9150509250925092565b5f60ff82169050919050565b610f6381610f4e565b82525050565b5f602082019050610f7c5f830184610f5a565b92915050565b5f60208284031215610f9757610f96610dd4565b5b5f610fa484828501610e1e565b91505092915050565b610fb681610df7565b82525050565b5f602082019050610fcf5f830184610fad565b92915050565b5f5f60408385031215610feb57610fea610dd4565b5b5f610ff885828601610e1e565b925050602061100985828601610e1e565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061105757607f821691505b60208210810361106a57611069611013565b5b50919050565b5f6060820190506110835f830186610fad565b6110906020830185610ed6565b61109d6040830184610ed6565b949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6110dc82610e32565b91506110e783610e32565b92508282019050808211156110ff576110fe6110a5565b5b9291505056fea264697066735822122013e328835bd3fa1f45865fa6a73996ca6799a72458b054d60e05498cf8150af164736f6c63430008210033a26469706673582212205fbc3dada0eaa713ae250e15dcdd04d51077e27e8fb419b667d6d5f01ecc922f64736f6c63430008210033",
}

// Dispenser is an auto generated Go binding around an Ethereum contract.
type Dispenser struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Dispenser) GetABI() abi.ABI {
	return c.abi
}

// NewDispenser creates a new instance of Dispenser.
func NewDispenser() *Dispenser {
	parsed, err := DispenserMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Dispenser{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Dispenser) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address owner) returns()
func (dispenser *Dispenser) PackConstructor(owner common.Address) []byte {
	enc, err := dispenser.abi.Pack("", owner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x85936228.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function assets(string ) view returns(address)
func (dispenser *Dispenser) PackAssets(arg0 string) []byte {
	enc, err := dispenser.abi.Pack("assets", arg0)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAssets is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x85936228.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function assets(string ) view returns(address)
func (dispenser *Dispenser) TryPackAssets(arg0 string) ([]byte, error) {
	return dispenser.abi.Pack("assets", arg0)
}

// UnpackAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x85936228.
//
// Solidity: function assets(string ) view returns(address)
func (dispenser *Dispenser) UnpackAssets(data []byte) (common.Address, error) {
	out, err := dispenser.abi.Unpack("assets", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackCreateAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xda7b7ce3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function createAsset(string name, string symbol, uint8 decimal) returns(address)
func (dispenser *Dispenser) PackCreateAsset(name string, symbol string, decimal uint8) []byte {
	enc, err := dispenser.abi.Pack("createAsset", name, symbol, decimal)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCreateAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xda7b7ce3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function createAsset(string name, string symbol, uint8 decimal) returns(address)
func (dispenser *Dispenser) TryPackCreateAsset(name string, symbol string, decimal uint8) ([]byte, error) {
	return dispenser.abi.Pack("createAsset", name, symbol, decimal)
}

// UnpackCreateAsset is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xda7b7ce3.
//
// Solidity: function createAsset(string name, string symbol, uint8 decimal) returns(address)
func (dispenser *Dispenser) UnpackCreateAsset(data []byte) (common.Address, error) {
	out, err := dispenser.abi.Unpack("createAsset", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackDispense is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x30568c34.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function dispense(string name, address to, uint256 amount) returns()
func (dispenser *Dispenser) PackDispense(name string, to common.Address, amount *big.Int) []byte {
	enc, err := dispenser.abi.Pack("dispense", name, to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDispense is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x30568c34.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function dispense(string name, address to, uint256 amount) returns()
func (dispenser *Dispenser) TryPackDispense(name string, to common.Address, amount *big.Int) ([]byte, error) {
	return dispenser.abi.Pack("dispense", name, to, amount)
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (dispenser *Dispenser) PackOwner() []byte {
	enc, err := dispenser.abi.Pack("owner")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function owner() view returns(address)
func (dispenser *Dispenser) TryPackOwner() ([]byte, error) {
	return dispenser.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (dispenser *Dispenser) UnpackOwner(data []byte) (common.Address, error) {
	out, err := dispenser.abi.Unpack("owner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceOwnership() returns()
func (dispenser *Dispenser) PackRenounceOwnership() []byte {
	enc, err := dispenser.abi.Pack("renounceOwnership")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x715018a6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceOwnership() returns()
func (dispenser *Dispenser) TryPackRenounceOwnership() ([]byte, error) {
	return dispenser.abi.Pack("renounceOwnership")
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (dispenser *Dispenser) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := dispenser.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (dispenser *Dispenser) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return dispenser.abi.Pack("transferOwnership", newOwner)
}

// DispenserOwnershipTransferred represents a OwnershipTransferred event raised by the Dispenser contract.
type DispenserOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const DispenserOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (DispenserOwnershipTransferred) ContractEventName() string {
	return DispenserOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (dispenser *Dispenser) UnpackOwnershipTransferredEvent(log *types.Log) (*DispenserOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != dispenser.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DispenserOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := dispenser.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range dispenser.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (dispenser *Dispenser) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], dispenser.abi.Errors["AssetAlreadyExists"].ID.Bytes()[:4]) {
		return dispenser.UnpackAssetAlreadyExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], dispenser.abi.Errors["AssetNotFound"].ID.Bytes()[:4]) {
		return dispenser.UnpackAssetNotFoundError(raw[4:])
	}
	if bytes.Equal(raw[:4], dispenser.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return dispenser.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], dispenser.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return dispenser.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// DispenserAssetAlreadyExists represents a AssetAlreadyExists error raised by the Dispenser contract.
type DispenserAssetAlreadyExists struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetAlreadyExists()
func DispenserAssetAlreadyExistsErrorID() common.Hash {
	return common.HexToHash("0xdc0d0aab66c63178623982491c96ac5a757e4c1592582973dd609e76a321a218")
}

// UnpackAssetAlreadyExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetAlreadyExists()
func (dispenser *Dispenser) UnpackAssetAlreadyExistsError(raw []byte) (*DispenserAssetAlreadyExists, error) {
	out := new(DispenserAssetAlreadyExists)
	if err := dispenser.abi.UnpackIntoInterface(out, "AssetAlreadyExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DispenserAssetNotFound represents a AssetNotFound error raised by the Dispenser contract.
type DispenserAssetNotFound struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetNotFound()
func DispenserAssetNotFoundErrorID() common.Hash {
	return common.HexToHash("0x470cbf473dda26c086d6854544465eb8cabe6ec3cb833bbf52751983b4af0f9f")
}

// UnpackAssetNotFoundError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetNotFound()
func (dispenser *Dispenser) UnpackAssetNotFoundError(raw []byte) (*DispenserAssetNotFound, error) {
	out := new(DispenserAssetNotFound)
	if err := dispenser.abi.UnpackIntoInterface(out, "AssetNotFound", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DispenserOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the Dispenser contract.
type DispenserOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func DispenserOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (dispenser *Dispenser) UnpackOwnableInvalidOwnerError(raw []byte) (*DispenserOwnableInvalidOwner, error) {
	out := new(DispenserOwnableInvalidOwner)
	if err := dispenser.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DispenserOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the Dispenser contract.
type DispenserOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func DispenserOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (dispenser *Dispenser) UnpackOwnableUnauthorizedAccountError(raw []byte) (*DispenserOwnableUnauthorizedAccount, error) {
	out := new(DispenserOwnableUnauthorizedAccount)
	if err := dispenser.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}
