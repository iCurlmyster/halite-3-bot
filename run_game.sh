#!/bin/sh

set -e
go build -o bot main

./halite --replay-directory replays/ -vvv --width 32 --height 32 "./bot" "./bot2"
