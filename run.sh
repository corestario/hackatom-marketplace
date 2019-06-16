#!/usr/bin/env bash

rm -rf ~/.hh*

make install

hhd init node_name --chain-id hhchain

hhcli keys add validator1 --recover <<< "12345678
base figure planet hazard sail easily honey advance tuition grab across unveil random kiss fence connect disagree evil recall latin cause brisk soft lunch
"

hhcli keys add alice --recover <<< "12345678
actor barely wait patrol moral amateur hole clerk misery truly salad wonder artefact orchard grit check abandon drip avoid shaft dirt thought melody drip
"

hhd add-genesis-account $(hhcli keys show validator1 -a) 1000token,100000000stake
hhd add-genesis-account $(hhcli keys show alice -a) 1000token

cp ./config.toml $HOME/.hhd/config/config.toml
cp ./config.toml $HOME/.hhcli/config/config.toml


hhcli config chain-id hhchain
hhcli config output json
hhcli config indent true
hhcli config trust-node true

hhd gentx --name validator1
hhd collect-gentxs
hhd validate-genesis

hhd start