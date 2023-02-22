# Initial Craton blockchain write in Go

## Introduction

This blockchain follows the classical Satoshi Paper with several adaptations to runs in private private way.

The blockchain receives the transactions with instructions of how to execute them, proofing it self when our outside partner says it is valid.

We divided the blockchain into three different parts:

- The blockchain, who stores the blocks and the transactions

- The wallet, who stores the private keys and the public keys.

- The gateway, who is the interface between the blockchains, the wallet and the outside world.

## Development

### First Way

The easiest way to get started is with docker-compose, clone the repository and start with:

```shell
docker-compose up
```

### Second Way

The second way is with Vagrant and Ansible, it will spin up a virtual machine inside your computer with the blockchain and the wallets and create a little network, it will simulate better the real world.

You need to have vagrant, ansible, docker and virtual-box installed to go in this path.
If you have all dependencies, just run the following command:

```shell
development_mode.sh
```

## Testing
To execute tests, specify the path of the test file: 

```shell
go test -v ./{path} 
```








