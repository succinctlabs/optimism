// SPDX-License-Identifier: MIT
pragma solidity >0.5.0;
pragma experimental ABIEncoderV2;

import './iOVM_L1ERC20Bridge.sol';

/**
 * @title iOVM_L1StandardBridge
 */
interface iOVM_L1StandardBridge is iOVM_L1ERC20Bridge {

    /**********
     * Events *
     **********/
    event ETHDepositInitiated(
        address indexed _from,
        address indexed _to,
        uint256 _amount,
        bytes _data
    );

    event ETHWithdrawalFinalized(
        address indexed _from,
        address indexed _to,
        uint256 _amount,
        bytes _data
    );

    /********************
     * Public Functions *
     ********************/

    function depositETH(
        uint32 _l2Gas,
        bytes calldata _data
    )
        external
        payable;

    function depositETHTo(
        address _to,
        uint32 _l2Gas,
        bytes calldata _data
    )
        external
        payable;

    /*************************
     * Cross-chain Functions *
     *************************/
    function finalizeETHWithdrawal(
        address _from,
        address _to,
        uint _amount,
        bytes calldata _data
    )
        external;
}
