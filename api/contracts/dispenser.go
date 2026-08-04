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

// ContractsMetaData contains all meta data concerning the Contracts contract.
var ContractsMetaData = bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"assets\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractPortmanAsset\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createAsset\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"precision\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dispense\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AssetAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AssetNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "Contracts",
	Bin: "0x608060405234801561000f575f5ffd5b506040516125de3803806125de833981810160405281019061003191906101d7565b805f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100a2575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100999190610211565b60405180910390fd5b6100b1816100b860201b60201c565b505061022a565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6101a68261017d565b9050919050565b6101b68161019c565b81146101c0575f5ffd5b50565b5f815190506101d1816101ad565b92915050565b5f602082840312156101ec576101eb610179565b5b5f6101f9848285016101c3565b91505092915050565b61020b8161019c565b82525050565b5f6020820190506102245f830184610202565b92915050565b6123a7806102375f395ff3fe608060405234801561000f575f5ffd5b5060043610610060575f3560e01c806330568c3414610064578063715018a614610080578063859362281461008a5780638da5cb5b146100ba578063da7b7ce3146100d8578063f2fde38b146100f4575b5f5ffd5b61007e600480360381019061007991906107ac565b610110565b005b61008861022c565b005b6100a4600480360381019061009f9190610818565b61023f565b6040516100b191906108ba565b60405180910390f35b6100c2610287565b6040516100cf91906108e2565b60405180910390f35b6100f260048036038101906100ed9190610931565b6102ae565b005b61010e600480360381019061010991906109b9565b6103f2565b005b610118610476565b5f6001846040516101299190610a36565b90815260200160405180910390205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036101be576040517f470cbf4700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff16637b1837de84846040518363ffffffff1660e01b81526004016101f9929190610a5b565b5f604051808303815f87803b158015610210575f5ffd5b505af1158015610222573d5f5f3e3d5ffd5b5050505050505050565b610234610476565b61023d5f6104fd565b565b6001818051602081018201805184825260208301602085012081835280955050505050505f915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6102b6610476565b5f73ffffffffffffffffffffffffffffffffffffffff166001846040516102dd9190610a36565b90815260200160405180910390205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610358576040517fdc0d0aab00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f838383604051610368906105c5565b61037493929190610ad9565b604051809103905ff08015801561038d573d5f5f3e3d5ffd5b509050806001856040516103a19190610a36565b90815260200160405180910390205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050505050565b6103fa610476565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361046a575f6040517f1e4fbdf700000000000000000000000000000000000000000000000000000000815260040161046191906108e2565b60405180910390fd5b610473816104fd565b50565b61047e6105be565b73ffffffffffffffffffffffffffffffffffffffff1661049c610287565b73ffffffffffffffffffffffffffffffffffffffff16146104fb576104bf6105be565b6040517f118cdaa70000000000000000000000000000000000000000000000000000000081526004016104f291906108e2565b60405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f33905090565b61185580610b1d83390190565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610631826105eb565b810181811067ffffffffffffffff821117156106505761064f6105fb565b5b80604052505050565b5f6106626105d2565b905061066e8282610628565b919050565b5f67ffffffffffffffff82111561068d5761068c6105fb565b5b610696826105eb565b9050602081019050919050565b828183375f83830152505050565b5f6106c36106be84610673565b610659565b9050828152602081018484840111156106df576106de6105e7565b5b6106ea8482856106a3565b509392505050565b5f82601f830112610706576107056105e3565b5b81356107168482602086016106b1565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6107488261071f565b9050919050565b6107588161073e565b8114610762575f5ffd5b50565b5f813590506107738161074f565b92915050565b5f819050919050565b61078b81610779565b8114610795575f5ffd5b50565b5f813590506107a681610782565b92915050565b5f5f5f606084860312156107c3576107c26105db565b5b5f84013567ffffffffffffffff8111156107e0576107df6105df565b5b6107ec868287016106f2565b93505060206107fd86828701610765565b925050604061080e86828701610798565b9150509250925092565b5f6020828403121561082d5761082c6105db565b5b5f82013567ffffffffffffffff81111561084a576108496105df565b5b610856848285016106f2565b91505092915050565b5f819050919050565b5f61088261087d6108788461071f565b61085f565b61071f565b9050919050565b5f61089382610868565b9050919050565b5f6108a482610889565b9050919050565b6108b48161089a565b82525050565b5f6020820190506108cd5f8301846108ab565b92915050565b6108dc8161073e565b82525050565b5f6020820190506108f55f8301846108d3565b92915050565b5f60ff82169050919050565b610910816108fb565b811461091a575f5ffd5b50565b5f8135905061092b81610907565b92915050565b5f5f5f60608486031215610948576109476105db565b5b5f84013567ffffffffffffffff811115610965576109646105df565b5b610971868287016106f2565b935050602084013567ffffffffffffffff811115610992576109916105df565b5b61099e868287016106f2565b92505060406109af8682870161091d565b9150509250925092565b5f602082840312156109ce576109cd6105db565b5b5f6109db84828501610765565b91505092915050565b5f81519050919050565b5f81905092915050565b8281835e5f83830152505050565b5f610a10826109e4565b610a1a81856109ee565b9350610a2a8185602086016109f8565b80840191505092915050565b5f610a418284610a06565b915081905092915050565b610a5581610779565b82525050565b5f604082019050610a6e5f8301856108d3565b610a7b6020830184610a4c565b9392505050565b5f82825260208201905092915050565b5f610a9c826109e4565b610aa68185610a82565b9350610ab68185602086016109f8565b610abf816105eb565b840191505092915050565b610ad3816108fb565b82525050565b5f6060820190508181035f830152610af18186610a92565b90508181036020830152610b058185610a92565b9050610b146040830184610aca565b94935050505056fe60a060405234801561000f575f5ffd5b5060405161185538038061185583398181016040528101906100319190610332565b338383816003908161004391906105db565b50806004908161005391906105db565b5050505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036100c6575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100bd91906106e9565b60405180910390fd5b6100d5816100ec60201b60201c565b508060ff1660808160ff1681525050505050610702565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61020e826101c8565b810181811067ffffffffffffffff8211171561022d5761022c6101d8565b5b80604052505050565b5f61023f6101af565b905061024b8282610205565b919050565b5f67ffffffffffffffff82111561026a576102696101d8565b5b610273826101c8565b9050602081019050919050565b8281835e5f83830152505050565b5f6102a061029b84610250565b610236565b9050828152602081018484840111156102bc576102bb6101c4565b5b6102c7848285610280565b509392505050565b5f82601f8301126102e3576102e26101c0565b5b81516102f384826020860161028e565b91505092915050565b5f60ff82169050919050565b610311816102fc565b811461031b575f5ffd5b50565b5f8151905061032c81610308565b92915050565b5f5f5f60608486031215610349576103486101b8565b5b5f84015167ffffffffffffffff811115610366576103656101bc565b5b610372868287016102cf565b935050602084015167ffffffffffffffff811115610393576103926101bc565b5b61039f868287016102cf565b92505060406103b08682870161031e565b9150509250925092565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061040857607f821691505b60208210810361041b5761041a6103c4565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261047d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610442565b6104878683610442565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6104cb6104c66104c18461049f565b6104a8565b61049f565b9050919050565b5f819050919050565b6104e4836104b1565b6104f86104f0826104d2565b84845461044e565b825550505050565b5f5f905090565b61050f610500565b61051a8184846104db565b505050565b5f5b82811015610540576105355f828401610507565b600181019050610521565b505050565b601f82111561059357828211156105925761055f81610421565b61056883610433565b61057185610433565b602086101561057e575f90505b80830161058d8284038261051f565b505050505b5b505050565b5f82821c905092915050565b5f6105b35f1984600802610598565b1980831691505092915050565b5f6105cb83836105a4565b9150826002028217905092915050565b6105e4826103ba565b67ffffffffffffffff8111156105fd576105fc6101d8565b5b61060782546103f1565b610612828285610545565b5f60209050601f831160018114610643575f8415610631578287015190505b61063b85826105c0565b8655506106a2565b601f19841661065186610421565b5f5b8281101561067857848901518255600182019150602085019450602081019050610653565b868310156106955784890151610691601f8916826105a4565b8355505b6001600288020188555050505b505050505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6106d3826106aa565b9050919050565b6106e3816106c9565b82525050565b5f6020820190506106fc5f8301846106da565b92915050565b60805161113b61071a5f395f610385015261113b5ff3fe608060405234801561000f575f5ffd5b50600436106100cd575f3560e01c8063715018a61161008a57806395d89b411161006457806395d89b41146101ff578063a9059cbb1461021d578063dd62ed3e1461024d578063f2fde38b1461027d576100cd565b8063715018a6146101bb5780637b1837de146101c55780638da5cb5b146101e1576100cd565b806306fdde03146100d1578063095ea7b3146100ef57806318160ddd1461011f57806323b872dd1461013d578063313ce5671461016d57806370a082311461018b575b5f5ffd5b6100d9610299565b6040516100e69190610db4565b60405180910390f35b61010960048036038101906101049190610e65565b610329565b6040516101169190610ebd565b60405180910390f35b61012761034b565b6040516101349190610ee5565b60405180910390f35b61015760048036038101906101529190610efe565b610354565b6040516101649190610ebd565b60405180910390f35b610175610382565b6040516101829190610f69565b60405180910390f35b6101a560048036038101906101a09190610f82565b6103a9565b6040516101b29190610ee5565b60405180910390f35b6101c36103ee565b005b6101df60048036038101906101da9190610e65565b610401565b005b6101e9610417565b6040516101f69190610fbc565b60405180910390f35b61020761043f565b6040516102149190610db4565b60405180910390f35b61023760048036038101906102329190610e65565b6104cf565b6040516102449190610ebd565b60405180910390f35b61026760048036038101906102629190610fd5565b6104f1565b6040516102749190610ee5565b60405180910390f35b61029760048036038101906102929190610f82565b610573565b005b6060600380546102a890611040565b80601f01602080910402602001604051908101604052809291908181526020018280546102d490611040565b801561031f5780601f106102f65761010080835404028352916020019161031f565b820191905f5260205f20905b81548152906001019060200180831161030257829003601f168201915b5050505050905090565b5f5f6103336105f7565b90506103408185856105fe565b600191505092915050565b5f600254905090565b5f5f61035e6105f7565b905061036b858285610610565b6103768585856106a3565b60019150509392505050565b5f7f0000000000000000000000000000000000000000000000000000000000000000905090565b5f5f5f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20549050919050565b6103f6610793565b6103ff5f61081a565b565b610409610793565b61041382826108dd565b5050565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60606004805461044e90611040565b80601f016020809104026020016040519081016040528092919081815260200182805461047a90611040565b80156104c55780601f1061049c576101008083540402835291602001916104c5565b820191905f5260205f20905b8154815290600101906020018083116104a857829003601f168201915b5050505050905090565b5f5f6104d96105f7565b90506104e68185856106a3565b600191505092915050565b5f60015f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905092915050565b61057b610793565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036105eb575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016105e29190610fbc565b60405180910390fd5b6105f48161081a565b50565b5f33905090565b61060b838383600161095c565b505050565b5f61061b84846104f1565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81101561069d578181101561068e578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161068593929190611070565b60405180910390fd5b61069c84848484035f61095c565b5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610713575f6040517f96c6fd1e00000000000000000000000000000000000000000000000000000000815260040161070a9190610fbc565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610783575f6040517fec442f0500000000000000000000000000000000000000000000000000000000815260040161077a9190610fbc565b60405180910390fd5b61078e838383610b2b565b505050565b61079b6105f7565b73ffffffffffffffffffffffffffffffffffffffff166107b9610417565b73ffffffffffffffffffffffffffffffffffffffff1614610818576107dc6105f7565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161080f9190610fbc565b60405180910390fd5b565b5f60055f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508160055f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361094d575f6040517fec442f050000000000000000000000000000000000000000000000000000000081526004016109449190610fbc565b60405180910390fd5b6109585f8383610b2b565b5050565b5f73ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036109cc575f6040517fe602df050000000000000000000000000000000000000000000000000000000081526004016109c39190610fbc565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610a3c575f6040517f94280d62000000000000000000000000000000000000000000000000000000008152600401610a339190610fbc565b60405180910390fd5b8160015f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508015610b25578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610b1c9190610ee5565b60405180910390a35b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1603610b7b578060025f828254610b6f91906110d2565b92505081905550610c49565b5f5f5f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610c04578381836040517fe450d38c000000000000000000000000000000000000000000000000000000008152600401610bfb93929190611070565b60405180910390fd5b8181035f5f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2081905550505b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610c90578060025f8282540392505081905550610cda565b805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610d379190610ee5565b60405180910390a3505050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f610d8682610d44565b610d908185610d4e565b9350610da0818560208601610d5e565b610da981610d6c565b840191505092915050565b5f6020820190508181035f830152610dcc8184610d7c565b905092915050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610e0182610dd8565b9050919050565b610e1181610df7565b8114610e1b575f5ffd5b50565b5f81359050610e2c81610e08565b92915050565b5f819050919050565b610e4481610e32565b8114610e4e575f5ffd5b50565b5f81359050610e5f81610e3b565b92915050565b5f5f60408385031215610e7b57610e7a610dd4565b5b5f610e8885828601610e1e565b9250506020610e9985828601610e51565b9150509250929050565b5f8115159050919050565b610eb781610ea3565b82525050565b5f602082019050610ed05f830184610eae565b92915050565b610edf81610e32565b82525050565b5f602082019050610ef85f830184610ed6565b92915050565b5f5f5f60608486031215610f1557610f14610dd4565b5b5f610f2286828701610e1e565b9350506020610f3386828701610e1e565b9250506040610f4486828701610e51565b9150509250925092565b5f60ff82169050919050565b610f6381610f4e565b82525050565b5f602082019050610f7c5f830184610f5a565b92915050565b5f60208284031215610f9757610f96610dd4565b5b5f610fa484828501610e1e565b91505092915050565b610fb681610df7565b82525050565b5f602082019050610fcf5f830184610fad565b92915050565b5f5f60408385031215610feb57610fea610dd4565b5b5f610ff885828601610e1e565b925050602061100985828601610e1e565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061105757607f821691505b60208210810361106a57611069611013565b5b50919050565b5f6060820190506110835f830186610fad565b6110906020830185610ed6565b61109d6040830184610ed6565b949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6110dc82610e32565b91506110e783610e32565b92508282019050808211156110ff576110fe6110a5565b5b9291505056fea26469706673582212200c5bfc85d4eada5bcc823ce5e7ed9f1fde9976601ef15cd20e9fbccc4a627c2a64736f6c63430008210033a264697066735822122014e9c0f1bf157229a97e06a534d08939f65b143aa24fe0c245c019c310eec90364736f6c63430008210033",
}

// Contracts is an auto generated Go binding around an Ethereum contract.
type Contracts struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *Contracts) GetABI() abi.ABI {
	return c.abi
}

// NewContracts creates a new instance of Contracts.
func NewContracts() *Contracts {
	parsed, err := ContractsMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &Contracts{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *Contracts) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address owner) returns()
func (contracts *Contracts) PackConstructor(owner common.Address) []byte {
	enc, err := contracts.abi.Pack("", owner)
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
func (contracts *Contracts) PackAssets(arg0 string) []byte {
	enc, err := contracts.abi.Pack("assets", arg0)
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
func (contracts *Contracts) TryPackAssets(arg0 string) ([]byte, error) {
	return contracts.abi.Pack("assets", arg0)
}

// UnpackAssets is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x85936228.
//
// Solidity: function assets(string ) view returns(address)
func (contracts *Contracts) UnpackAssets(data []byte) (common.Address, error) {
	out, err := contracts.abi.Unpack("assets", data)
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
// Solidity: function createAsset(string name, string symbol, uint8 precision) returns()
func (contracts *Contracts) PackCreateAsset(name string, symbol string, precision uint8) []byte {
	enc, err := contracts.abi.Pack("createAsset", name, symbol, precision)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCreateAsset is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xda7b7ce3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function createAsset(string name, string symbol, uint8 precision) returns()
func (contracts *Contracts) TryPackCreateAsset(name string, symbol string, precision uint8) ([]byte, error) {
	return contracts.abi.Pack("createAsset", name, symbol, precision)
}

// PackDispense is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x30568c34.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function dispense(string name, address to, uint256 amount) returns()
func (contracts *Contracts) PackDispense(name string, to common.Address, amount *big.Int) []byte {
	enc, err := contracts.abi.Pack("dispense", name, to, amount)
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
func (contracts *Contracts) TryPackDispense(name string, to common.Address, amount *big.Int) ([]byte, error) {
	return contracts.abi.Pack("dispense", name, to, amount)
}

// PackOwner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8da5cb5b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function owner() view returns(address)
func (contracts *Contracts) PackOwner() []byte {
	enc, err := contracts.abi.Pack("owner")
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
func (contracts *Contracts) TryPackOwner() ([]byte, error) {
	return contracts.abi.Pack("owner")
}

// UnpackOwner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (contracts *Contracts) UnpackOwner(data []byte) (common.Address, error) {
	out, err := contracts.abi.Unpack("owner", data)
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
func (contracts *Contracts) PackRenounceOwnership() []byte {
	enc, err := contracts.abi.Pack("renounceOwnership")
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
func (contracts *Contracts) TryPackRenounceOwnership() ([]byte, error) {
	return contracts.abi.Pack("renounceOwnership")
}

// PackTransferOwnership is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf2fde38b.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (contracts *Contracts) PackTransferOwnership(newOwner common.Address) []byte {
	enc, err := contracts.abi.Pack("transferOwnership", newOwner)
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
func (contracts *Contracts) TryPackTransferOwnership(newOwner common.Address) ([]byte, error) {
	return contracts.abi.Pack("transferOwnership", newOwner)
}

// ContractsOwnershipTransferred represents a OwnershipTransferred event raised by the Contracts contract.
type ContractsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           *types.Log // Blockchain specific contextual infos
}

const ContractsOwnershipTransferredEventName = "OwnershipTransferred"

// ContractEventName returns the user-defined event name.
func (ContractsOwnershipTransferred) ContractEventName() string {
	return ContractsOwnershipTransferredEventName
}

// UnpackOwnershipTransferredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (contracts *Contracts) UnpackOwnershipTransferredEvent(log *types.Log) (*ContractsOwnershipTransferred, error) {
	event := "OwnershipTransferred"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != contracts.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(ContractsOwnershipTransferred)
	if len(log.Data) > 0 {
		if err := contracts.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contracts.abi.Events[event].Inputs {
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
func (contracts *Contracts) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], contracts.abi.Errors["AssetAlreadyExists"].ID.Bytes()[:4]) {
		return contracts.UnpackAssetAlreadyExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], contracts.abi.Errors["AssetNotFound"].ID.Bytes()[:4]) {
		return contracts.UnpackAssetNotFoundError(raw[4:])
	}
	if bytes.Equal(raw[:4], contracts.abi.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]) {
		return contracts.UnpackOwnableInvalidOwnerError(raw[4:])
	}
	if bytes.Equal(raw[:4], contracts.abi.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]) {
		return contracts.UnpackOwnableUnauthorizedAccountError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// ContractsAssetAlreadyExists represents a AssetAlreadyExists error raised by the Contracts contract.
type ContractsAssetAlreadyExists struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetAlreadyExists()
func ContractsAssetAlreadyExistsErrorID() common.Hash {
	return common.HexToHash("0xdc0d0aab66c63178623982491c96ac5a757e4c1592582973dd609e76a321a218")
}

// UnpackAssetAlreadyExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetAlreadyExists()
func (contracts *Contracts) UnpackAssetAlreadyExistsError(raw []byte) (*ContractsAssetAlreadyExists, error) {
	out := new(ContractsAssetAlreadyExists)
	if err := contracts.abi.UnpackIntoInterface(out, "AssetAlreadyExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractsAssetNotFound represents a AssetNotFound error raised by the Contracts contract.
type ContractsAssetNotFound struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AssetNotFound()
func ContractsAssetNotFoundErrorID() common.Hash {
	return common.HexToHash("0x470cbf473dda26c086d6854544465eb8cabe6ec3cb833bbf52751983b4af0f9f")
}

// UnpackAssetNotFoundError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AssetNotFound()
func (contracts *Contracts) UnpackAssetNotFoundError(raw []byte) (*ContractsAssetNotFound, error) {
	out := new(ContractsAssetNotFound)
	if err := contracts.abi.UnpackIntoInterface(out, "AssetNotFound", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractsOwnableInvalidOwner represents a OwnableInvalidOwner error raised by the Contracts contract.
type ContractsOwnableInvalidOwner struct {
	Owner common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableInvalidOwner(address owner)
func ContractsOwnableInvalidOwnerErrorID() common.Hash {
	return common.HexToHash("0x1e4fbdf7f3ef8bcaa855599e3abf48b232380f183f08f6f813d9ffa5bd585188")
}

// UnpackOwnableInvalidOwnerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableInvalidOwner(address owner)
func (contracts *Contracts) UnpackOwnableInvalidOwnerError(raw []byte) (*ContractsOwnableInvalidOwner, error) {
	out := new(ContractsOwnableInvalidOwner)
	if err := contracts.abi.UnpackIntoInterface(out, "OwnableInvalidOwner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// ContractsOwnableUnauthorizedAccount represents a OwnableUnauthorizedAccount error raised by the Contracts contract.
type ContractsOwnableUnauthorizedAccount struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func ContractsOwnableUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0x118cdaa7a341953d1887a2245fd6665d741c67c8c50581daa59e1d03373fa188")
}

// UnpackOwnableUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OwnableUnauthorizedAccount(address account)
func (contracts *Contracts) UnpackOwnableUnauthorizedAccountError(raw []byte) (*ContractsOwnableUnauthorizedAccount, error) {
	out := new(ContractsOwnableUnauthorizedAccount)
	if err := contracts.abi.UnpackIntoInterface(out, "OwnableUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}
