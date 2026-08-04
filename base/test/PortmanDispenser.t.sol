// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Test} from "forge-std/Test.sol";
import {Ownable} from "openzeppelin-contracts/contracts/access/Ownable.sol";
import {PortmanAsset, PortmanDispenser} from "../src/PortmanDispenser.sol";

contract PortmanDispenserTest is Test {
    PortmanDispenser public dispenser;
    address public admin;
    address public alice;
    address public bob;

    function setUp() public {
        admin = vm.addr(1);
        alice = vm.addr(2);
        bob = vm.addr(3);

        dispenser = new PortmanDispenser(admin);
    }

    function testAssetCreation() public {
        string memory name = "TestToken";
        string memory symbol = "TTK";
        uint8 precision = 18;

        // Admin can create an asset
        vm.prank(admin);
        dispenser.createAsset(name, symbol, precision);
        PortmanAsset asset = dispenser.assets(name);

        assertFalse(address(asset) == address(0), "Asset should exist");
        assertEq(asset.name(), name, "Asset name should match");
        assertEq(asset.symbol(), symbol, "Asset symbol should match");
        assertEq(asset.decimals(), precision, "Asset precision should match");

        // Non-admin cannot create an asset
        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, alice));
        dispenser.createAsset(name, symbol, precision);

        // Admin cannot create duplicate asset
        vm.prank(admin);
        vm.expectRevert(abi.encodeWithSelector(PortmanDispenser.AssetAlreadyExists.selector));
        dispenser.createAsset(name, symbol, precision);

        // Admin cannot directly make an asset fund an account
        vm.prank(admin);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, admin));
        asset.fund(alice, 1000 * (10 ** precision));
    }

    function testDispense() public {
        string memory name = "TestToken";
        string memory symbol = "TTK";
        uint8 precision = 18;
        uint256 amount = 1000 * (10 ** precision);

        vm.prank(admin);
        dispenser.createAsset(name, symbol, precision);
        PortmanAsset asset = dispenser.assets(name);

        vm.prank(admin);
        dispenser.dispense(name, alice, amount);
        assertEq(asset.balanceOf(alice), amount, "Alice should receive the correct amount");

        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, bob));
        dispenser.dispense(name, bob, amount);
    }
}
