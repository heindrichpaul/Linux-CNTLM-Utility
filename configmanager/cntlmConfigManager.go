package configmanager

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/heindrichpaul/Linux-CNTLM-Utility/credentialmanager"
	"github.com/heindrichpaul/Linux-CNTLM-Utility/domain"
	"golang.org/x/sys/unix"
)

func replaceIfNeeded(line string, credentials *credentialmanager.Credentials) (newline string) {

	if strings.Contains(line, "Username") && !strings.Contains(strings.Fields(line)[1], credentials.Username) {
		newline = fmt.Sprintf("Username    %s", credentials.Username)
	} else if strings.Contains(line, "Domain") && !strings.Contains(strings.Fields(line)[1], credentials.Domain) {
		newline = fmt.Sprintf("Domain    %s", credentials.Domain)
	} else if strings.Contains(line, "PassLM") && !credentials.UseClearTextPassword {
		newline = credentials.PasswordHashes[0]
	} else if strings.Contains(line, "PassLM") && credentials.UseClearTextPassword {
		newline = "#" + line
	} else if strings.Contains(line, "PassNTLMv2") && !credentials.UseClearTextPassword {
		newline = credentials.PasswordHashes[2]
	} else if strings.Contains(line, "PassNTLMv2") && credentials.UseClearTextPassword {
		newline = "#" + line
	} else if strings.Contains(line, "PassNT") && credentials.UseClearTextPassword {
		newline = "#" + line
	} else if strings.Contains(line, "PassNT") && !credentials.UseClearTextPassword {
		newline = credentials.PasswordHashes[1]
	} else if strings.Contains(line, "Password") && !credentials.UseClearTextPassword {
		newline = fmt.Sprintf("#Password CLEARTEXT PASSWORD HAS BEEN REMOVED")
	} else if strings.Contains(line, "Password") && credentials.UseClearTextPassword {
		newline = fmt.Sprintf("Password    %s", credentials.Password)
	} else {
		newline = line
	}

	return
}

func UpdateCntlmConfig(config *domain.CntlmConfig, credentials *credentialmanager.Credentials) error {
	err := updateCntlmConfig(config.CntlmConfigPath, credentials)
	if err != nil {
		return err
	}
	if config.Profiles != nil {
		for _, profile := range config.Profiles {
			err = updateCntlmConfig(profile.ProfileFileLocation, credentials)
			return err
		}
	}
	return err
}

func updateCntlmConfig(filename string, credentials *credentialmanager.Credentials) error {

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		lines[i] = replaceIfNeeded(line, credentials)
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func isFileReadyToBeWrittenTo(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("The file %s does not exist", path)
	}

	access := unix.Access(path, unix.W_OK) == nil

	if !access {
		return fmt.Errorf("this process does not have access to the file %s. Maybe run the command with sudo if it needs root permissions", path)
	}

	return nil
}

func SetUpUtility(config *domain.CntlmConfig) (err error) {

	reader := bufio.NewReader(os.Stdin)

	useClearTextPassword := false
	filename := ""

	if len(config.CntlmConfigPath) == 0 {
		fmt.Printf("Enter path to cntlm config file (/etc/cntlm.conf): ")
		filename, err = reader.ReadString('\n')
		if err != nil {
			return err
		}

		filename = strings.TrimSpace(filename)

		if len(filename) < 1 {
			filename = "/etc/cntlm.conf"
		}
		config.CntlmConfigPath = filename
	}

	if config.PasswordProperties == nil {

		err = isFileReadyToBeWrittenTo(filename)

		fmt.Printf("Did you want to save your password in clear text in the configuration file (y/N): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		response = strings.TrimSpace(response)

		if strings.EqualFold(response, "Y") {
			useClearTextPassword = true
		}
		config.PasswordProperties = &domain.Password{}
		config.PasswordProperties.UseClearTextPassword = useClearTextPassword
	}

	return
}
