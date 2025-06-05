import Web3 from 'web3';

let web3;

if (window.ethereum) {
    web3 = new Web3(window.ethereum);
    window.ethereum.enable();
} else {
    alert("Please install MetaMask to use this app");
}

export default web3;
