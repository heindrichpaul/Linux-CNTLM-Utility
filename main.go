package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

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

	err = config.SaveConfig()
	if err != nil {
		log.Fatalln(err)
	}

	utility := utilitymanager.NewUtilityManger(config)

	menu(config, utility)

}

func readConfig() *domain.CntlmConfig {
	configPath := os.Getenv("CNTLM_UTILITY_CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = fmt.Sprintf("%s/.cntlm/config.json", os.Getenv("HOME"))
	}

	config, _ := domain.LoadJSON(configPath)
	return config
}

func menu(config *domain.CntlmConfig, utility *utilitymanager.UtilityManager) {
	exit := false
	reader := bufio.NewReader(os.Stdin)

	for !exit {
		printBanner()
		fmt.Printf("Please select an action to perform:\n1) Change credentials\n2) Switch Profiles\n3) Create Profile\n4) Exit\n")
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		temp := strings.Split(text, "\n")
		if len(temp) == 2 {
			choice, err := strconv.Atoi(temp[0])
			if err == nil {
				switch choice {
				case 1:
					changeCredentials(config, utility)
				case 2:
					switchProfiles(config, utility)
				case 3:
					createProfile(utility)
				case 4:
					exit = true
				default:

				}
				clearConsole()
			}
		}

	}
}

func changeCredentials(config *domain.CntlmConfig, utility *utilitymanager.UtilityManager) {
	domain, username, password := credentialmanager.GetCredentials()

	credentials := credentialmanager.CreateCredentials(username, domain, password, config.PasswordProperties.UseClearTextPassword)

	err := configmanager.UpdateCntlmConfig(config, credentials)
	if err != nil {
		log.Fatalln("An error occured during cntlm config update")
	}

	err = utility.RestartCntlm()
	if err != nil {
		log.Fatalln(err)
	}
}

func switchProfiles(config *domain.CntlmConfig, utility *utilitymanager.UtilityManager) {
	profiles := make(map[int]string, 0)
	for i, profile := range config.Profiles {
		profiles[i] = profile.Name
	}

	for id, profile := range profiles {
		fmt.Printf("%d) %s", id, profile)
	}
	fmt.Print("Please select the profile you wish to switch to: ")
	reader := bufio.NewReader(os.Stdin)
	profileNumber, _ := reader.ReadString('\n')
	temp := strings.Split(profileNumber, "\n")
	if len(temp) == 2 {
		choice, err := strconv.Atoi(temp[0])
		if err == nil {
			profileName, ok := profiles[choice]
			if ok {
				utility.SwitchProfile(profileName)
				err := utility.RestartCntlm()
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}

}

func createProfile(utility *utilitymanager.UtilityManager) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the name of the profile:  ")
	profileName, _ := reader.ReadString('\n')
	fmt.Print("Please enter the location of the profile config:  ")
	profileLocation, _ := reader.ReadString('\n')
	utility.CreateProfile(profileName, profileLocation)
}

func clearConsole() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		log.Fatalf("Unable to clear the terminal window.")
	}
}

func printBanner() {
	fmt.Println(`===============================================================================`)
	fmt.Println(`  _    _                 ___ _  _ _____ _    __  __   _   _ _   _ _ _ _        `)
	fmt.Println(` | |  (_)_ _ _  ___ __  / __| \| |_   _| |  |  \/  | | | | | |_(_) (_) |_ _  _ `)
	fmt.Println(` | |__| | ' \ || \ \ / | (__|  " | | | | |__| |\/| | | |_| |  _| | | |  _| || |`)
	fmt.Println(` |____|_|_||_\_,_/_\_\  \___|_|\_| |_| |____|_|  |_|  \___/ \__|_|_|_|\__|\_, |`)
	fmt.Println(`                                                                          |__/ `)
	fmt.Println(`===============================================================================`)
}
