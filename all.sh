#!/bin/bash
echo "Start..."
echo `date`
echo
echo
./raspberryCompile.sh
./linuxComp.sh
echo
echo
echo `date`
echo "...end"
echo "ok... running it"
./rssnest
cp ./data/newData.json ./web/public/newData.json
