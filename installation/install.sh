#!/bin/sh
# Written by Heindrich Friedrich Paul
# Except explicitly stated, all right reserved by Heindrich Friedrich Paul and 
# the script is licensed under the MIT license.


#Moves executable to bin directory for execution on the path
sudo mv changecntlmpassword /bin/changecntlmpassword
#Changes ownership to root
sudo chown root:root /bin/changecntlmpassword
#Changes permissions to add execution flag
sudo chmod 755 /bin/changecntlmpassword
#Creates config directory in user home
mkdir -p $HOME/.cntlm/
#Moves default config into place
mv config.json $HOME/.cntlm/config.json
#Changes ownership of config to root
sudo chown root:root $HOME/.cntlm/config.json
#Changes file permissions of config file
sudo chmod 644 $HOME/.cntlm/config.json
