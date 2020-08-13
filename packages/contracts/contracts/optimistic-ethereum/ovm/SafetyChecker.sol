pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

/* Library Imports */
import { ContractResolver } from "../utils/resolvers/ContractResolver.sol";

import { console } from "@nomiclabs/buidler/console.sol";

/**
 * @title SafetyChecker
 * @notice Safety Checker contract used to check whether or not bytecode is
 *         safe, meaning:
 *              1. It uses only whitelisted opcodes.
 *              2. All CALLs are to the Execution Manager and have no value.
 */
contract SafetyChecker is ContractResolver {
    /*
     * Constructor
     */

    /**
     * @param _addressResolver Address of the AddressResolver contract.
     */
    constructor(
        address _addressResolver
    )
        public
        ContractResolver(_addressResolver)
    {
    }


    /*
     * Public Functions
     */

    /**
     * Returns whether or not all of the provided bytecode is safe.
     * @dev More info on creation vs. runtime bytecode:
     * https://medium.com/authereum/bytecode-and-init-code-and-runtime-code-oh-my-7bcd89065904.
     * @param _bytecode The bytecode to safety check. This can be either
     *                  creation bytecode (aka initcode) or runtime bytecode
     *                  (aka cont
     * More info on creation vs. runtime bytecode:
     * https://medium.com/authereum/bytecode-and-init-code-and-runtime-code-oh-my-7bcd89065904ract code).
     * @return `true` if the bytecode is safe, `false` otherwise.
     */
    function isBytecodeSafe(
        bytes memory _bytecode
    )
        public
        view
        returns (bool)
    {
        // autogenerated by gen_safety_checker_constants.py
        /*uint256[8] memory skip = [
          uint256(0x0001010101010101010101010000000001010101010101010101010101010000),
          uint256(0x0100000000000000000000000000000000000000010101010101000000010100),
          uint256(0x0000000000000000000000000000000001010101000000010101010100000000),
          uint256(0x0203040500000000000000000000000000000000000000000000000000000000),
          uint256(0x0101010101010101010101010101010101010101010101010101010101010101),
          uint256(0x0101010101000000000000000000000000000000000000000000000000000000),
          uint256(0x0000000000000000000000000000000000000000000000000000000000000000),
          uint256(0x0000000000000000000000000000000000000000000000000000000000000000)];*/
        uint256 _opcodeProcMask = ~uint256(0xffffffffffffffffffffffe000000000fffffffff070ffff9c0ffffec000f001);
        uint256 _opcodeStopMask = ~uint256(0x6008000000000000000000000000000000000000004000000000000000000001);
        uint256 _opcodePushMask = ~uint256(0xffffffff000000000000000000000000);

        uint256 codeLength;
        uint256 _pc;
        assembly {
            _pc := add(_bytecode, 0x20)
        }
        codeLength = _pc + _bytecode.length;
        do {
            // current opcode: 0x00...0xff
            uint256 op;

            // inline assembly removes the extra add + bounds check
            assembly {
                let tmp := mload(_pc)

                // this works, it just isn't fast
                let mpc := 0

                // this is fast
                /*mpc := add(mpc, byte(0, mload(add(skip, byte(mpc, tmp)))))
                mpc := add(mpc, byte(0, mload(add(skip, byte(mpc, tmp)))))
                mpc := add(mpc, byte(0, mload(add(skip, byte(mpc, tmp)))))
                mpc := add(mpc, byte(0, mload(add(skip, byte(mpc, tmp)))))
                mpc := add(mpc, byte(0, mload(add(skip, byte(mpc, tmp)))))
                mpc := add(mpc, byte(0, mload(add(skip, byte(mpc, tmp)))))*/

                // footer
                _pc := add(_pc, mpc)
                op := byte(mpc, tmp)
            }

            // + push opcodes
            // + stop opcodes [STOP(0x00),JUMP(0x56),RETURN(0xf3),REVERT(0xfd),INVALID(0xfe)]
            // + caller opcode CALLER(0x33) (which is technically on the blacklist)
            // + blacklisted opcodes
            uint256 opBit = 1 << op;
            if (opBit & _opcodeProcMask == 0) {
                if (opBit & _opcodePushMask == 0) {
                    // subsequent bytes are not opcodes. Skip them.
                    _pc += (op - 0x5e);
                    // all pushes are valid opcodes
                    continue;
                } else if (opBit & _opcodeStopMask == 0) {
                    // STOP or JUMP or RETURN or REVERT or INVALID (see safety checker docs in wiki for more info)
                    // We are now inside unreachable code until we hit a JUMPDEST!
                    do {
                        _pc++;
                        assembly {
                            op := byte(0, mload(_pc))
                        }
                        if (op == 0x5b) break;
                        if ((1 << op) & _opcodePushMask == 0) _pc += (op - 0x5f);
                    } while (_pc < codeLength);
                    // op is 0x5b, so we don't continue here since the _pc++ is fine
                } else if (op == 0x33) {
                    // Sequence around CALLER must be:
                    // 1. CALLER (execution manager address) <-- We are here
                    // 2. PUSH1 0x0
                    // 3. SWAP1
                    // 4. GAS (gas for call)
                    // 5. CALL

                    uint256 ops;
                    assembly {
                        ops := shr(208, mload(_pc))
                    }

                    // allowed = CALLER PUSH1 0x0 SWAP1 GAS CALL
                    if (ops != 0x336000905af1) {
                        console.log('Encountered a bad call');
                        return false;
                    }

                    _pc += 6;
                    continue;
                } else {
                    // encountered a non-whitelisted opcode!
                    console.log('Encountered a non-whitelisted opcode (in decimal):', op, "at location", _pc);
                    return false;
                }
            }
            _pc++;
        } while (_pc < codeLength);
        return true;
    }
}
