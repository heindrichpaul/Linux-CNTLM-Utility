package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Credentials struct {
	Username             string
	Password             string
	PasswordHashes       []string
	Domain               string
	useClearTextPassword bool
}

func main() {

	filename, useClearTextPassword, err := setUpUtility()
	if err != nil {
		log.Fatalln(err)
	}

	domain, username, password := getCredentials()

	credentials := createCredentials(username, domain, password, useClearTextPassword)

	err = updateCntlmConfig(filename, credentials)
	if err != nil {
		log.Fatalln("An error occured during cntlm config update")
	}

	err = restartCntlm()
	if err != nil {
		log.Fatalln(err)
	}

}

func setUpUtility() (filename string, useClearTextPassword bool, err error) {

	reader := bufio.NewReader(os.Stdin)

	useClearTextPassword = false

	fmt.Printf("Enter path to cntlm config file (/etc/cntlm.conf): ")
	filename, err = reader.ReadString('\n')
	if err != nil {
		return "", false, err
	}

	filename = strings.TrimSpace(filename)

	if len(filename) < 1 {
		filename = "/etc/cntlm.conf"
	}

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

	return
}

func getCredentials() (domain string, username string, password string) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	if strings.ContainsRune(username, '/') {
		fmt.Printf("Did you not mean a back-slash instead of a forward-slash? Should we fix it? (Y/n): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		response = strings.TrimSpace(response)

		if len(response) == 0 || strings.EqualFold(response, "Y") {
			username = strings.Replace(username, "/", "\\", 1)
		}

	}

	if !strings.Contains(username, "\\") {
		fmt.Printf("Enter Domain: ")
		domain, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		credentials := strings.Split(username, "\\")
		if len(credentials) != 2 {
			log.Fatalln("Could not split the domain and username")
		}
		domain = credentials[0]
		username = credentials[1]

	}

	fmt.Printf("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\n")

	domain = strings.TrimSpace(domain)
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(string(bytePassword))

	return
}

func createCntlmHashes(domain, username, password string) (hashes []string, err error) {

	cmdName := "cntlm"
	cmdArgs := []string{"-H", "-d", domain, "-u", username, "-p", password}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln("Error creating StdoutPipe for Cmd", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			output := scanner.Text()
			if strings.Contains(output, "Password:") {
				continue
			} else {
				hashes = append(hashes, output)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatalln("Error starting Cmd", err)
	}

	err = cmd.Wait()
	if err != nil {

		log.Fatalln("Error waiting Cmd", err)
	}

	return hashes, nil
}

func updateCntlmConfig(filename string, credentials *Credentials) error {

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

func restartCntlm() error {

	if os.Geteuid() == 0 {

		cmdName := "systemctl"
		cmdArgs := []string{"restart", "cntlm"}

		cmd := exec.Command(cmdName, cmdArgs...)

		fmt.Printf("Restarting the cntlm service\n")

		err := cmd.Start()
		if err != nil {
			log.Fatalln("Error restarting cntlm", err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Fatalln("Error waiting for the cntlm cmd", err)
		}

		fmt.Printf("Done restarting the cntlm service\n")

	} else {
		log.Printf("You have to run the command to restart the cntlm service\n")
	}

	return nil
}

func isFileReadyToBeWrittenTo(path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("The file %s does not exist", path)
	}

	access := unix.Access(path, unix.W_OK) == nil

	if !access {
		return fmt.Errorf("This process does not have access to the file %s. Maybe run the command with sudo if it needs root permissions.", path)
	}

	return nil
}

func createCredentials(username, domain, password string, useClearTextPassword bool) *Credentials {

	credentials := Credentials{}
	credentials.Username = username
	credentials.Domain = domain
	credentials.useClearTextPassword = useClearTextPassword
	if useClearTextPassword {
		credentials.Password = password
	} else {
		hashes, err := createCntlmHashes(domain, username, password)
		if err != nil {
			log.Fatalln(err)
		}
		credentials.PasswordHashes = hashes
	}
	return &credentials
}

func replaceIfNeeded(line string, credentials *Credentials) (newline string) {

	if strings.Contains(line, "Username") && !strings.Contains(strings.Fields(line)[1], credentials.Username) {
		newline = fmt.Sprintf("Username    %s", credentials.Username)
	} else if strings.Contains(line, "Domain") && !strings.Contains(strings.Fields(line)[1], credentials.Domain) {
		newline = fmt.Sprintf("Domain    %s", credentials.Domain)
	} else if strings.Contains(line, "PassLM") && credentials.useClearTextPassword == false {
		newline = credentials.PasswordHashes[0]
	} else if strings.Contains(line, "PassLM") && credentials.useClearTextPassword == true {
		newline = "#" + line
	} else if strings.Contains(line, "PassNTLMv2") && credentials.useClearTextPassword == false {
		newline = credentials.PasswordHashes[2]
	} else if strings.Contains(line, "PassNTLMv2") && credentials.useClearTextPassword == true {
		newline = "#" + line
	} else if strings.Contains(line, "PassNT") && credentials.useClearTextPassword == true {
		newline = "#" + line
	} else if strings.Contains(line, "PassNT") && credentials.useClearTextPassword == false {
		newline = credentials.PasswordHashes[1]
	} else if strings.Contains(line, "Password") && credentials.useClearTextPassword == false {
		newline = fmt.Sprintf("#Password CLEARTEXT PASSWORD HAS BEEN REMOVED")
	} else if strings.Contains(line, "Password") && credentials.useClearTextPassword == true {
		newline = fmt.Sprintf("Password    %s", credentials.Password)
	} else {
		newline = line
	}

	return
}
