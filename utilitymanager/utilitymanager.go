package utilitymanager

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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
