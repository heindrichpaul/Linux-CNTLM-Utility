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

func main() {
	filename, err := getCntlmConfig()
	if err != nil {
		log.Fatalln(err)
	}

	domain, username, password := credentials()

	hashes, err := getCntlmHashes(domain, username, password)
	if err != nil {
		log.Fatalln(err)
	}

	err = updateCntlmConfig(filename, domain, username, hashes)
	if err != nil {
		log.Fatalln("An error occured during cntlm config update")
	}

	err = restartCntlm()
	if err != nil {
		log.Fatalln(err)
	}
}

func credentials() (domain string, username string, password string) {
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

func getCntlmHashes(domain, username, password string) (hashes []string, err error) {
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
			hashes = append(hashes, scanner.Text())
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

	if len(hashes) != 4 {
		return nil, fmt.Errorf("There was an error getting the complete output from the Cntlm command")
	}

	hashes = hashes[1:]

	return hashes, nil
}

func updateCntlmConfig(filename, domain, username string, hashes []string) error {

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "Username") && !strings.Contains(strings.Fields(line)[1], username) {
			lines[i] = fmt.Sprintf("Username    %s", username)
		} else if strings.Contains(line, "Domain") && !strings.Contains(strings.Fields(line)[1], domain) {
			lines[i] = fmt.Sprintf("Domain    %s", domain)
		} else if strings.Contains(line, "PassLM") {
			lines[i] = hashes[0]
		} else if strings.Contains(line, "PassNTLMv2") {
			lines[i] = hashes[2]
		} else if strings.Contains(line, "PassNT") {
			lines[i] = hashes[1]
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func getCntlmConfig() (filename string, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter path to cntlm config file: ")
	filename, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	filename = strings.TrimSpace(filename)

	err = writable(filename)

	return

}

func writable(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("The file %s does not exist", path)
	}

	access := unix.Access(path, unix.W_OK) == nil

	if !access {
		return fmt.Errorf("This process does not have access to the file %s. Maybe run the command with sudo if it needs root permissions.", path)
	}

	return nil
}

func restartCntlm() error {
	cmdName := "systemctl"
	cmdArgs := []string{"restart", "cntlm"}

	cmd := exec.Command(cmdName, cmdArgs...)

	err := cmd.Start()
	if err != nil {
		log.Fatalln("Error restarting cntlm", err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalln("Error waiting Cmd", err)
	}

	return nil
}
