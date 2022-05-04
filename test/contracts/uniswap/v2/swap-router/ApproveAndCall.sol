// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity =0.7.6;
pragma abicoder v2;

import './IERC20.sol';
import './INonfungiblePositionManager.sol';

import './IApproveAndCall.sol';
import './ImmutableState.sol';

/// @title Approve and Call
/// @notice Allows callers to approve the Uniswap V3 position manager from this contract,
/// for any token, and then make calls into the position manager
abstract contract ApproveAndCall is IApproveAndCall, ImmutableState {
    function tryApprove(address token, uint256 amount) private returns (bool) {
        (bool success, bytes memory data) =
            token.call(abi.encodeWithSelector(IERC20.approve.selector, positionManager, amount));
        return success && (data.length == 0 || abi.decode(data, (bool)));
    }

    /// @inheritdoc IApproveAndCall
    function getApprovalType(address token, uint256 amount) external override returns (ApprovalType) {
        // check existing approval
        if (IERC20(token).allowance(address(this), positionManager) >= amount) return ApprovalType.NOT_REQUIRED;

        // try type(uint256).max / type(uint256).max - 1
        if (tryApprove(token, type(uint256).max)) return ApprovalType.MAX;
        if (tryApprove(token, type(uint256).max - 1)) return ApprovalType.MAX_MINUS_ONE;

        // set approval to 0 (must succeed)
        require(tryApprove(token, 0));

        // try type(uint256).max / type(uint256).max - 1
        if (tryApprove(token, type(uint256).max)) return ApprovalType.ZERO_THEN_MAX;
        if (tryApprove(token, type(uint256).max - 1)) return ApprovalType.ZERO_THEN_MAX_MINUS_ONE;

        revert();
    }

    /// @inheritdoc IApproveAndCall
    function approveMax(address token) external payable override {
        require(tryApprove(token, type(uint256).max));
    }

    /// @inheritdoc IApproveAndCall
    function approveMaxMinusOne(address token) external payable override {
        require(tryApprove(token, type(uint256).max - 1));
    }

    /// @inheritdoc IApproveAndCall
    function approveZeroThenMax(address token) external payable override {
        require(tryApprove(token, 0));
        require(tryApprove(token, type(uint256).max));
    }

    /// @inheritdoc IApproveAndCall
    function approveZeroThenMaxMinusOne(address token) external payable override {
        require(tryApprove(token, 0));
        require(tryApprove(token, type(uint256).max - 1));
    }

    /// @inheritdoc IApproveAndCall
    function callPositionManager(bytes memory data) public payable override returns (bytes memory result) {
        bool success;
        (success, result) = positionManager.call(data);

        if (!success) {
            // Next 5 lines from https://ethereum.stackexchange.com/a/83577
            if (result.length < 68) revert();
            assembly {
                result := add(result, 0x04)
            }
            revert(abi.decode(result, (string)));
        }
    }

    function balanceOf(address token) private view returns (uint256) {
        return IERC20(token).balanceOf(address(this));
    }

    /// @inheritdoc IApproveAndCall
    function mint(MintParams calldata params) external payable override returns (bytes memory result) {
        return
            callPositionManager(
                abi.encodeWithSelector(
                    INonfungiblePositionManager.mint.selector,
                    INonfungiblePositionManager.MintParams({
                        token0: params.token0,
                        token1: params.token1,
                        fee: params.fee,
                        tickLower: params.tickLower,
                        tickUpper: params.tickUpper,
                        amount0Desired: balanceOf(params.token0),
                        amount1Desired: balanceOf(params.token1),
                        amount0Min: params.amount0Min,
                        amount1Min: params.amount1Min,
                        recipient: params.recipient,
                        deadline: type(uint256).max // deadline should be checked via multicall
                    })
                )
            );
    }

    /// @inheritdoc IApproveAndCall
    function increaseLiquidity(IncreaseLiquidityParams calldata params)
        external
        payable
        override
        returns (bytes memory result)
    {
        return
            callPositionManager(
                abi.encodeWithSelector(
                    INonfungiblePositionManager.increaseLiquidity.selector,
                    INonfungiblePositionManager.IncreaseLiquidityParams({
                        tokenId: params.tokenId,
                        amount0Desired: balanceOf(params.token0),
                        amount1Desired: balanceOf(params.token1),
                        amount0Min: params.amount0Min,
                        amount1Min: params.amount1Min,
                        deadline: type(uint256).max // deadline should be checked via multicall
                    })
                )
            );
    }
}
