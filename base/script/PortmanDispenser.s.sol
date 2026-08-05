// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {PortmanDispenser} from "../src/PortmanDispenser.sol";

contract PortmanDispenserScript is Script {
    function run() public {
        address sender = vm.envAddress("SENDER");

        vm.startBroadcast(sender);
        PortmanDispenser dispenser = new PortmanDispenser(sender);
        console.log("PortmanDispenser:", address(dispenser));

        address asset = dispenser.createAsset("Portman Dollar", "USD", 6);
        console.log("PortmanAsset:", asset);
        vm.stopBroadcast();
    }
}
