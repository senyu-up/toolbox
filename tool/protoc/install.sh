#!/bin/zsh
unzip -d $HOME/.local/protoc protoc-3.14.0-osx-x86_64.zip
echo "PATH=\$PATH:$HOME/.local/protoc/bin" >> $HOME/.zshrc
source "$HOME"/.zshrc
protoc --version