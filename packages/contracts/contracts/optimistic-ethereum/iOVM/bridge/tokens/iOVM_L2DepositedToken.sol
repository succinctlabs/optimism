// SPDX-License-Identifier: MIT
pragma solidity >0.5.0;
pragma experimental ABIEncoderV2;

/**
 * @title iOVM_L2DepositedToken
 */
interface iOVM_L2DepositedToken {

    /**********
     * Events *
     **********/

    event WithdrawalInitiated(
        address indexed _from,
        address indexed _to,
        uint256 _amount
    );

    event DepositFinalized(
        address indexed _from,
        address indexed _to,
        uint256 _amount
    );


    /********************
     * Public Functions *
     ********************/

    function withdraw(
        uint _amount,
        bytes calldata _data
    )
        external;

    function withdrawTo(
        address _to,
        uint _amount,
        bytes calldata _data
    )
        external;


    /*************************
     * Cross-chain Functions *
     *************************/

    function finalizeDeposit(
        address _from,
        address _to,
        uint _amount,
        bytes calldata _data
    )
        external;
}
