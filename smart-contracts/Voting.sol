// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Voting {
    struct Candidate {
        string name;
        uint voteCount;
    }

    address public owner;
    mapping(address => bool) public voters;
    Candidate[] public candidates;

    constructor(string[] memory candidateNames) {
        owner = msg.sender;
        for (uint i = 0; i < candidateNames.length; i++) {
            candidates.push(Candidate({name: candidateNames[i], voteCount: 0}));
        }
    }

    function vote(uint candidateIndex) public {
        require(!voters[msg.sender], "Already voted.");
        require(candidateIndex < candidates.length, "Invalid candidate.");

        voters[msg.sender] = true;
        candidates[candidateIndex].voteCount += 1;
    }

    function getCandidates() public view returns (Candidate[] memory) {
        return candidates;
    }
}
