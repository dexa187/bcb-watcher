pragma solidity ^0.5.0;

contract ERC20 {
    event TransferWithData(address indexed from, address indexed to, uint256 tokens, bytes data);
    event Approval(address indexed tokenOwner, address indexed spender, uint tokens);
}