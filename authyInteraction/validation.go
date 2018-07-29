package authyInteraction

import (
	"fmt"
	"os"
)

func init() {
	neededEnvVars := []string{"Authy_API_KEY"}
	for i := 0; i < len(neededEnvVars); i++ {
		if os.Getenv(neededEnvVars[i]) == "" {
			fmt.Printf("\nERROR: env variable %s not found\n", neededEnvVars[i])
		}
	}
}
