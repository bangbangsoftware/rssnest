RSSNEST
=======

A little program to download audio/visual content from a list of rss feeds.Also tweet and create a web page from the content, as well as the gold and silver prices.

Target computer is a raspberry pi, seems to have memory leaks already, should help sharpen my golang skills.


To compile for the Raspberry pi
-------------------------------
Have to set some env variables for compile for the pi...

* export GOARM=5
* export GOOS=linux
* export GOARCH=arm
* export GOPATH=/home/mick/work/rssnest

