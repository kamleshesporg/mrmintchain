#!/bin/bash

echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./mrmintchaind validate-genesis --home /mrmintchain

echo "starting mrmintchain node $ID in background ..."
./mrmintchaind start \
--home /mrmintchain \
--keyring-backend test

echo "started mrmintchain node"
tail -f /dev/null