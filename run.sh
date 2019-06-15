#!/usr/bin/env bash

rm -rf ~/.hh*

make install

hhd init node_name --chain-id hhchain

hhcli keys add validator1

hhd add-genesis-account $(hhcli keys show validator1 -a) 1000hhtoken,100000000stake

hhcli config chain-id hhchain
hhcli config output json
hhcli config indent true
hhcli config trust-node true

hhd gentx --name validator1
hhd collect-gentxs
hhd validate-genesis

hhd start