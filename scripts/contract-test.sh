#!/bin/bash

KEY="mykey"
CHAINID="mrmintchain_9000-1"
MONIKER="localtestnet"

# stop and remove existing daemon and client data and process(es)
rm -rf ~/.mrmintchain*
pkill -f "mrmintchain*"

make build-mrmintchain

# if $KEY exists it should be override
"$PWD"/build/mrmintchaind keys add $KEY --keyring-backend test --algo "eth_secp256k1"

# Set moniker and chain-id for mrmintchain (Moniker can be anything, chain-id must be an integer)
"$PWD"/build/mrmintchaind init $MONIKER --chain-id $CHAINID

# Change parameter token denominations to mnt
cat $HOME/.mrmintchain/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="mnt"' > $HOME/.mrmintchain/config/tmp_genesis.json && mv $HOME/.mrmintchain/config/tmp_genesis.json $HOME/.mrmintchain/config/genesis.json
cat $HOME/.mrmintchain/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="mnt"' > $HOME/.mrmintchain/config/tmp_genesis.json && mv $HOME/.mrmintchain/config/tmp_genesis.json $HOME/.mrmintchain/config/genesis.json
cat $HOME/.mrmintchain/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="mnt"' > $HOME/.mrmintchain/config/tmp_genesis.json && mv $HOME/.mrmintchain/config/tmp_genesis.json $HOME/.mrmintchain/config/genesis.json
cat $HOME/.mrmintchain/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="mnt"' > $HOME/.mrmintchain/config/tmp_genesis.json && mv $HOME/.mrmintchain/config/tmp_genesis.json $HOME/.mrmintchain/config/genesis.json

# Allocate genesis accounts (cosmos formatted addresses)
"$PWD"/build/mrmintchaind add-genesis-account "$("$PWD"/build/mrmintchaind keys show "$KEY" -a --keyring-backend test)" 100000000000000000000mnt,10000000000000000000mnt --keyring-backend test

# Sign genesis transaction
"$PWD"/build/mrmintchaind gentx $KEY 10000000000000000000mnt --amount=100000000000000000000mnt --keyring-backend test --chain-id $CHAINID

# Collect genesis tx
"$PWD"/build/mrmintchaind collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
"$PWD"/build/mrmintchaind validate-genesis

# Start the node (remove the --pruning=nothing flag if historical queries are not needed) in background and log to file
"$PWD"/build/mrmintchaind start --pruning=nothing --rpc.unsafe --json-rpc.address="127.0.0.1:8545" --keyring-backend test > mrmintchaind.log 2>&1 &

# Give mrmintchaind node enough time to launch
sleep 5

solcjs --abi "$PWD"/tests/solidity/suites/basic/contracts/Counter.sol --bin -o "$PWD"/tests/solidity/suites/basic/counter
mv "$PWD"/tests/solidity/suites/basic/counter/*.abi "$PWD"/tests/solidity/suites/basic/counter/counter_sol.abi 2> /dev/null
mv "$PWD"/tests/solidity/suites/basic/counter/*.bin "$PWD"/tests/solidity/suites/basic/counter/counter_sol.bin 2> /dev/null

# Query for the account
ACCT=$(curl --fail --silent -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' -H "Content-Type: application/json" http://localhost:8545 | grep -o '\0x[^"]*')
echo "$ACCT"

# Start testcases (not supported)
# curl -X POST --data '{"jsonrpc":"2.0","method":"personal_unlockAccount","params":["'$ACCT'", ""],"id":1}' -H "Content-Type: application/json" http://localhost:8545

#PRIVKEY="$("$PWD"/build/mrmintchaind keys export $KEY)"

## need to get the private key from the account in order to check this functionality.
cd tests/solidity/suites/basic/ && go get && go run main.go $ACCT

# After tests
# kill test mrmintchaind
echo "going to shutdown mrmintchaind in 3 seconds..."
sleep 3
pkill -f "mrmintchain*"