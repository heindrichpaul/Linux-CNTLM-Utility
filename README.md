# Linux CNTLM Utility[![Go Report Card](https://goreportcard.com/badge/github.com/heindrichpaul/Linux-CNTLM-Utility)](https://goreportcard.com/report/github.com/heindrichpaul/Linux-CNTLM-Utility)

This utility can be used in a scenario to simplify the process of changing a user's password in the user's CNTLM's config file where corporations require the regular change of a user's password.


## Table of Content
* [Documentation](#documentation)
  * [Tool Requirements](#requirements)
  * [Use](#use)  
    *[Change credentials](#change-credentials)  
    *[Switch Profiles](#switch-profiles)  
    *[Create Profile](#create-profile)  
  * [Config](#config)  
### <a name="documentaion">Documentation
#### <a name="requirements"></a>Tool Requirements
This utility was written to support the systemd linux distributions.

#### <a name="use"></a>Use
The utility was rewritten to make the process even easier. For this a menu was created as well as support for multiple profiles.

The main menu when executing the utility is displayed below:

```
===============================================================================
  _    _                 ___ _  _ _____ _    __  __   _   _ _   _ _ _ _        
 | |  (_)_ _ _  ___ __  / __| \| |_   _| |  |  \/  | | | | | |_(_) (_) |_ _  _ 
 | |__| | ' \ || \ \ / | (__|  " | | | | |__| |\/| | | |_| |  _| | | |  _| || |
 |____|_|_||_\_,_/_\_\  \___|_|\_| |_| |____|_|  |_|  \___/ \__|_|_|_|\__|\_, |
                                                                          |__/ 
===============================================================================
Please select an action to perform:
1) Change credentials
2) Switch Profiles
3) Create Profile
4) Exit
Enter text: 
```

In the menu above you have three main options ([Change credentials](#change-credentials), [Switch Profiles](#switch-profiles), [Create Profile](#create-profile)) and an exit option.

##### <a name="change-credentials"></a>Change credentials
When selecting "Change credentials" the following prompts will appear:

```
Enter Username: 
Enter Domain: 
Enter New Password: 
Confirm Password:
```

These prompts work exactly like version 1
##### <a name="switch-profiles"></a>Switch Profiles
When selecting "Switch Profiles" the following menu will appear:
```
0) Home
1) Office
Please select the profile you wish to switch to:
```

The default config file already contains a Home and Office place holder. If you wish to edit them, you will have to edit the config file. For more clarity on how to edit the config file, read the config section.
##### <a name="create-profile"></a>Create Profile
When selecting "Create Profile" the prompts below will appear to assist you in creating a new profile. This can also be done manually by editing the config file.

```
Please enter the name of the profile:  demo
Please enter the location of the profile config:  /etc/cntlm.conf_demo
```

When creating a profile replace ```demo``` with the name of your profile and ```/etc/cntlm.conf_demo``` with the path to the copy of the ```etc/cntlm.conf``` file that contains the modifications for the profile.

It is important to note that the file referenced by the property ```profileFileLocation``` should be a valid cntlm config file with your extra modifications for the profile. When switching profiles the app will copy the profile config over the existing one.

#### <a name="config"></a>Config
If you wish to edit the file you will have to edit it at the default location ``` $HOME\.cntlm\config.json ``` with an editor with root permissions. The default location of the config file can be overridden by using the environmental variable ```CNTLM_UTILITY_CONFIG_PATH```.

The default config file is shown below:
```
{
    "cntlmConfigPath":"/etc/cntlm.conf",
    "passwordProperties":
        {
            "useClearTextPassword":false
        },
    "profiles":[
        {
            "name":"Home",
            "profileFileLocation":""
        },
        {
            "name":"Office",
            "profileFileLocation":""
        }
    ]
}
```
