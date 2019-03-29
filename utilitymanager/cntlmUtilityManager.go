package utilitymanager

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/heindrichpaul/Linux-CNTLM-Utility/domain"
)

type UtilityManager struct {
	config *domain.CntlmConfig
}

func NewUtilityManger(config *domain.CntlmConfig) *UtilityManager {
	z := &UtilityManager{
		config: config,
	}

	return z
}

func (z *UtilityManager) SwitchProfile(name string) {
	if z.config.Profiles != nil {
		for _, profile := range z.config.Profiles {
			if strings.EqualFold(profile.Name, name) {
				err := copy(profile.ProfileFileLocation, z.config.CntlmConfigPath)
				if err != nil {
					log.Fatalf("An error occured while switching profiles: %s\n", err.Error())
				}
				return
			}
		}
		log.Printf("No profile with name %s was found.\n", name)
	}
}

func (z *UtilityManager) RestartCntlm() error {

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
		fmt.Printf("You have to run the command to restart the cntlm service\n")
	}

	return nil
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func (z *UtilityManager) CreateProfile(name, profileLocationPath string) {
	profile := domain.Profile{
		Name:                name,
		ProfileFileLocation: profileLocationPath,
	}

	z.config.Profiles = append(z.config.Profiles, profile)
	err := z.config.SaveConfig()
	if err != nil {
		log.Println("Error in saving the config file.")
	}
}
