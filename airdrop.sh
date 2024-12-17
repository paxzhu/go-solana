#!/bin/bash
PUBKEY=$(cat wallet_pubkey.txt)
solana airdrop 1 $PUBKEY
