package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func UnmarshalCntlmConfig(data []byte) (CntlmConfig, error) {
	var r CntlmConfig
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CntlmConfig) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CntlmConfig struct {
	CntlmConfigPath    string    `json:"cntlmConfigPath"`
	PasswordProperties *Password `json:"passwordProperties"`
	Profiles           []Profile `json:"profiles"`
}

type Password struct {
	UseClearTextPassword bool `json:"useClearTextPassword"`
}

type Profile struct {
	Name                string `json:"name"`
	ProfileFileLocation string `json:"profileFileLocation"`
}

func LoadJSON(path string) (*CntlmConfig, error) {
	var config CntlmConfig

	file, err := os.Open(path)
	if err != nil {
		return &config, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return &config, err
	}

	return &config, err
}

func (config *CntlmConfig) SaveConfig() error {
	configPath := os.Getenv("CNTLM_UTILITY_CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = fmt.Sprintf("%s/.cntlm/config.json", os.Getenv("HOME"))
	}

	err := SaveJSON(configPath, config)
	if err != nil {
		log.Fatalf("an error occured while saving the config: %s", err.Error())
	}
	return err
}

func SaveJSON(path string, config *CntlmConfig) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	return err
}
