// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title InsurancePool
 * @dev Manages spike insurance purchases and payouts (Upgradeable)
 */
contract InsurancePool is Initializable, OwnableUpgradeable, ReentrancyGuardUpgradeable {
    IERC20 public usdt;
    
    // Insurance policy parameters
    uint256 public premiumAmount;
    uint256 public coverageAmount;
    uint256 public coverageDuration;
    
    struct Policy {
        address user;
        uint256 premium;
        uint256 coverageAmount;
        uint256 purchaseTime;
        uint256 expiryTime;
        bool active;
        bool claimed;
    }
    
    // Mapping from user address to their policies
    mapping(address => Policy[]) public userPolicies;
    
    // Backend oracle address authorized to trigger payouts
    address public oracle;
    
    // Events
    event PolicyPurchased(address indexed user, uint256 policyId, uint256 premium, uint256 coverage, uint256 expiryTime);
    event PayoutExecuted(address indexed user, uint256 policyId, uint256 spikeId, uint256 amount);
    event OracleUpdated(address indexed newOracle);
    
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    /**
     * @dev Initialize the contract (replaces constructor for upgradeable contracts)
     */
    function initialize(address _usdt) public initializer {
        __Ownable_init(msg.sender);
        __ReentrancyGuard_init();
        
        usdt = IERC20(_usdt);
        oracle = msg.sender;
        premiumAmount = 10 * 10**6; // 10 USDT (6 decimals)
        coverageAmount = 100 * 10**6; // 100 USDT payout
        coverageDuration = 24 hours; // Policy valid for 24 hours
    }
    
    /**
     * @dev Purchase insurance policy
     */
    function buyInsurance() external nonReentrant {
        // Transfer premium from user
        require(usdt.transferFrom(msg.sender, address(this), premiumAmount), "Premium transfer failed");
        
        // Create new policy
        Policy memory newPolicy = Policy({
            user: msg.sender,
            premium: premiumAmount,
            coverageAmount: coverageAmount,
            purchaseTime: block.timestamp,
            expiryTime: block.timestamp + coverageDuration,
            active: true,
            claimed: false
        });
        
        userPolicies[msg.sender].push(newPolicy);
        uint256 policyId = userPolicies[msg.sender].length - 1;
        
        emit PolicyPurchased(msg.sender, policyId, premiumAmount, coverageAmount, newPolicy.expiryTime);
    }
    
    /**
     * @dev Execute payout when spike is detected (called by backend oracle)
     */
    function executePayout(address user, uint256 policyId, uint256 spikeId) external nonReentrant {
        require(msg.sender == oracle, "Only oracle can execute payout");
        require(policyId < userPolicies[user].length, "Invalid policy ID");
        
        Policy storage policy = userPolicies[user][policyId];
        
        require(policy.active, "Policy not active");
        require(!policy.claimed, "Already claimed");
        require(block.timestamp <= policy.expiryTime, "Policy expired");
        require(usdt.balanceOf(address(this)) >= policy.coverageAmount, "Insufficient pool balance");
        
        // Mark as claimed
        policy.claimed = true;
        policy.active = false;
        
        // Transfer payout
        require(usdt.transfer(user, policy.coverageAmount), "Payout transfer failed");
        
        emit PayoutExecuted(user, policyId, spikeId, policy.coverageAmount);
    }
    
    /**
     * @dev Get user's active policies count
     */
    function getUserPoliciesCount(address user) external view returns (uint256) {
        return userPolicies[user].length;
    }
    
    /**
     * @dev Get specific policy details
     */
    function getPolicy(address user, uint256 policyId) external view returns (Policy memory) {
        require(policyId < userPolicies[user].length, "Invalid policy ID");
        return userPolicies[user][policyId];
    }
    
    /**
     * @dev Check if user has active unclaimed policy
     */
    function hasActivePolicy(address user) external view returns (bool) {
        Policy[] memory policies = userPolicies[user];
        for (uint256 i = 0; i < policies.length; i++) {
            if (policies[i].active && !policies[i].claimed && block.timestamp <= policies[i].expiryTime) {
                return true;
            }
        }
        return false;
    }
    
    /**
     * @dev Get all policies for a specific user
     * @param user Address of the user
     * @return Array of all policies for the user
     */
    function getUserPolicies(address user) external view returns (Policy[] memory) {
        return userPolicies[user];
    }
    
    /**
     * @dev Update oracle address
     */
    function setOracle(address _oracle) external onlyOwner {
        oracle = _oracle;
        emit OracleUpdated(_oracle);
    }
    
    /**
     * @dev Update policy parameters
     */
    function setPolicyParams(uint256 _premium, uint256 _coverage, uint256 _duration) external onlyOwner {
        premiumAmount = _premium;
        coverageAmount = _coverage;
        coverageDuration = _duration;
    }
    
    /**
     * @dev Fund the pool (for demo purposes)
     */
    function fundPool(uint256 amount) external {
        require(usdt.transferFrom(msg.sender, address(this), amount), "Fund transfer failed");
    }
    
    /**
     * @dev Get pool balance
     */
    function getPoolBalance() external view returns (uint256) {
        return usdt.balanceOf(address(this));
    }
    /**
     * @dev Clear all user policies for the provided list of users (only owner)
     * Owner must provide the list of all users who have policies.
     */
    function clearAllUserPolicies(address[] calldata users) external onlyOwner {
        for (uint256 i = 0; i < users.length; i++) {
            delete userPolicies[users[i]];
        }
    }
}