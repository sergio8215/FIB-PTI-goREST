#!/bin/bash

ps aux | grep apt
sudo kill PROCNUM
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
