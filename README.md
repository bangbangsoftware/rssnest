RSSNEST
=======

A little program to download audio/visual content from a list of rss feeds. Also tweet and create a web page from the content, as well as the gold and silver prices.

Target computer is a raspberry pi.

Quick start 
===========

GO
--

Assuming you are using linux there is a handy compile script linuxComp.sh, it will need changing to have your go path in. Run that and it will leave you with the rssnest executable. Then run

./rssnest -conf example_conf.json

This will no doubt download a copy of linux voice podcast into directory . and create some json, then crash as it cannot ftp the data onto the webserver. To get it working as intended:

*  The example_conf.json will need replacing with real values, 
*  As will the casts json referenced in the conf.json


JS
--
I guess this is not so much a quick start as you'll need to set all the web stuff up as well

* sudo npm install jspm -g
* cd ./web
* jspm install
* cd public

And ftp these files on to the webserver keeping the directory structure:

* ./jspm_packages/github/google/material-design-lite@1.0.4/material.min.css
* ./jspm_packages/system.js
* ./config.js
* ./js/main.js


That should do it...

Stop the growth
===============
There two problems with file growth on this:

* rssnest.log      - This will just keep growing
* alreadyHave.json - Also this will keep growing and is a quick and dirty way to persist some data, without setting up a database etc.

TODO
====

* Stop the growth
* Refactor
* Write some tests
* Write a quick go program to help ftp the assets to webserver?
