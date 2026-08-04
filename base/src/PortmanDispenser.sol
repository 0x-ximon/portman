// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Ownable} from "openzeppelin-contracts/contracts/access/Ownable.sol";
import {ERC20} from "openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract PortmanAsset is ERC20, Ownable {
    uint8 private immutable DECIMALS;

    constructor(string memory name, string memory symbol, uint8 _decimals) ERC20(name, symbol) Ownable(msg.sender) {
        DECIMALS = _decimals;
    }

    function fund(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }

    function decimals() public view override returns (uint8) {
        return DECIMALS;
    }
}

contract PortmanDispenser is Ownable {
    mapping(string => PortmanAsset) public assets;

    error AssetAlreadyExists();
    error AssetNotFound();

    constructor(address owner) Ownable(owner) {}

    function dispense(string memory name, address to, uint256 amount) external onlyOwner {
        PortmanAsset asset = assets[name];
        require(address(asset) != address(0), AssetNotFound());
        asset.fund(to, amount);
    }

    function createAsset(string memory name, string memory symbol, uint8 precision) external onlyOwner {
        require(address(assets[name]) == address(0), AssetAlreadyExists());
        PortmanAsset asset = new PortmanAsset(name, symbol, precision);
        assets[name] = asset;
    }
}

