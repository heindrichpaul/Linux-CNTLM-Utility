#!/bin/sh
# Written by Heindrich Friedrich Paul
# Except explicitly stated, all right reserved by Heindrich Friedrich Paul and 
# the script is licensed under the MIT license.


#Moves executable to bin directory for execution on the path
mv changecntlmpassword $HOME/bin
sudo chown $USER:$USER $HOME/bin/changecntlmpassword
chmod 755 $HOME/bin/changecntlmpassword
mkdir -p $HOME/.cntlm/
mv config.json $HOME/.cntlm/config.json
sudo chown $USER:$USER $HOME/.cntlm/config.json
chmod 644 $HOME/.cntlm/config.json
