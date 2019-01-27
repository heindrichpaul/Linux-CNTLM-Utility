package credentialmanager

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Credentials struct {
	Username             string
	Password             string
	PasswordHashes       []string
	Domain               string
	UseClearTextPassword bool
}

func GetCredentials() (domain string, username string, password string) {

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

	passwordConfirmed := false
	for passwordConfirmed == false {

		fmt.Printf("Enter New Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("\n")

		fmt.Printf("Confirm Password: ")
		bytePassword2, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("\n")

		password = strings.TrimSpace(string(bytePassword))
		password2 := strings.TrimSpace(string(bytePassword2))

		if password == password2 {
			passwordConfirmed = true
		} else {
			fmt.Printf("The passwords did not match please enter it again.\n")
		}
	}

	domain = strings.TrimSpace(domain)
	username = strings.TrimSpace(username)

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

func CreateCredentials(username, domain, password string, useClearTextPassword bool) *Credentials {

	credentials := Credentials{}
	credentials.Username = username
	credentials.Domain = domain
	credentials.UseClearTextPassword = useClearTextPassword
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
