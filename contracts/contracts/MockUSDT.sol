// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
 * @title MockUSDT
 * @dev Simple ERC20 token for testnet demo purposes
 */
contract MockUSDT is ERC20 {
    constructor() ERC20("Mock USDT", "mUSDT") {
        // Mint 1 million tokens to deployer for testing
        _mint(msg.sender, 1000000 * 10**decimals());
    }

    /**
     * @dev Allow anyone to mint tokens for testing
     */
    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }

    /**
     * @dev Override decimals to match USDT (6 decimals)
     */
    function decimals() public pure override returns (uint8) {
        return 6;
    }
}
