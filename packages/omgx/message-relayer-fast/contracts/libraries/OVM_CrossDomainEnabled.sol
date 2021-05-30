// SPDX-License-Identifier: MIT
pragma solidity >0.5.0 <0.8.0;
/* Interface Imports */
import { iAbs_BaseCrossDomainMessenger } from "@eth-optimism/contracts/iOVM/bridge/messaging/iAbs_BaseCrossDomainMessenger.sol";

/**
 * @title OVM_CrossDomainEnabled
 * @dev Helper contract for contracts performing cross-domain communications
 *
 * Compiler used: defined by inheriting contract
 * Runtime target: defined by inheriting contract
 */
contract OVM_CrossDomainEnabled {
    // Messenger contract used to send and recieve messages from the other domain.
    address public senderMessenger;
    address public relayerMessenger;

    /***************
     * Constructor *
     ***************/    
    constructor(
        address _senderMessenger,
        address _relayerMessenger
    ) {
        senderMessenger = _senderMessenger;
        relayerMessenger = _relayerMessenger;
    }

    /**********************
     * Function Modifiers *
     **********************/

    /**
     * @notice Enforces that the modified function is only callable by a specific cross-domain account.
     * @param _sourceDomainAccount The only account on the originating domain which is authenticated to call this function.
     */
    modifier onlyFromCrossDomainAccount(
        address _sourceDomainAccount
    ) {
        require(
            msg.sender == address(getCrossDomainRelayerMessenger()),
            "OVM_XCHAIN: messenger contract unauthenticated"
        );

        require(
            getCrossDomainRelayerMessenger().xDomainMessageSender() == _sourceDomainAccount,
            "OVM_XCHAIN: wrong sender of cross-domain message"
        );

        _;
    }
    
    /**********************
     * Internal Functions *
     **********************/

    /**
     * @notice Gets the messenger, usually from storage.  This function is exposed in case a child contract needs to override.
     * @return The address of the cross-domain messenger contract which should be used. 
     */
    function getCrossDomainSenderMessenger()
        internal
        virtual
        returns(
            iAbs_BaseCrossDomainMessenger
        )
    {
        return iAbs_BaseCrossDomainMessenger(senderMessenger);
    }

    /**
     * @notice Gets the messenger, usually from storage.  This function is exposed in case a child contract needs to override.
     * @return The address of the cross-domain messenger contract which should be used. 
     */
    function getCrossDomainRelayerMessenger()
        internal
        virtual
        returns(
            iAbs_BaseCrossDomainMessenger
        )
    {
        return iAbs_BaseCrossDomainMessenger(relayerMessenger);
    }

    /**
     * @notice Sends a message to an account on another domain
     * @param _crossDomainTarget The intended recipient on the destination domain
     * @param _data The data to send to the target (usually calldata to a function with `onlyFromCrossDomainAccount()`)
     * @param _gasLimit The gasLimit for the receipt of the message on the target domain.
     */
    function sendCrossDomainMessage(
        address _crossDomainTarget,
        bytes memory _data,
        uint32 _gasLimit
    ) internal {
        getCrossDomainSenderMessenger().sendMessage(_crossDomainTarget, _data, _gasLimit);
    }
}
