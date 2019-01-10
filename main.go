package main

import (
	"fmt"
	"log"
	"os"

	"github.com/heindrichpaul/Linux-CNTLM-Utility/domain"

	"github.com/heindrichpaul/Linux-CNTLM-Utility/configmanager"
	"github.com/heindrichpaul/Linux-CNTLM-Utility/credentialmanager"
	"github.com/heindrichpaul/Linux-CNTLM-Utility/utilitymanager"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	config := readConfig()
	err := configmanager.SetUpUtility(config)
	if err != nil {
		log.Fatalln(err)
	}

	err = storeConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	domain, username, password := credentialmanager.GetCredentials()

	credentials := credentialmanager.CreateCredentials(username, domain, password, config.PasswordProperties.UseClearTextPassword)

	err = configmanager.UpdateCntlmConfig(config.CntlmConfigPath, credentials)
	if err != nil {
		log.Fatalln("An error occured during cntlm config update")
	}

	utility := utilitymanager.NewUtilityManger(config)

	err = utility.RestartCntlm()
	if err != nil {
		log.Fatalln(err)
	}

}

func readConfig() *domain.CntlmConfig {
	configPath := os.Getenv("CNTLM_UTILITY_CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = fmt.Sprintf("%s/.cntlm/config.json", os.Getenv("HOME"))
	}

	config, _ := domain.LoadJSON(configPath)
	/*if err != nil {
		log.Fatalf("an error occured while loading the config: %s", err.Error())
	}*/
	return config
}

func storeConfig(config *domain.CntlmConfig) error {
	configPath := os.Getenv("CNTLM_UTILITY_CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = fmt.Sprintf("%s/.cntlm/config.json", os.Getenv("HOME"))
	}

	err := domain.SaveJSON(configPath, config)
	if err != nil {
		log.Fatalf("an error occured while saving the config: %s", err.Error())
	}
	return err
}
func menu() {

}
