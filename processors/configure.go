package processors

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/nanobox-io/nanobox/models"
)

var configured bool

func Configure() error {
	// make sure to only run configure one time
	if configured {
		return nil
	}
	configured = true

	<-time.After(time.Second)

	config := &models.Config{
		Provider:  "docker-machine",
		MountType: "native",
		CPUs:      1,
		RAM:       1,
	}

	fmt.Print(`
CONFIGURE NANOBOX
---------------------------------------------------------------
Please answer the following questions so we can customize your
nanobox configuration. Feel free to update your config at any
time by running: 'nanobox configure'

(Learn more at : https://docs.nanobox.io/local-config/nanobox-config-yml/)
`)

	// ask about provider
	config.Provider = stringAsker(`
How would you like to run nanobox?
  a) Inside a lightweight VM
  b) Via Docker Native

(recommended a)
Answer: `, map[string]string{"a": "docker-machine", "b": "native"})

	// if provider == docker-machine ask more questions
	if config.Provider == "native" {
		config.Save()
		return nil
	}

	// ask about cpus
	config.CPUs = intAsker(fmt.Sprintf(`
How many CPU cores would you like to make available to the VM (1-%d)?

(recommended > 2)
Answer: `, runtime.NumCPU()), runtime.NumCPU())

	// ask about ram
	config.RAM = intAsker(`
How many GB of RAM would you like to make available to the VM (2-4)?

(recommended > 1)
Answer: `, 4)

	// ask about mount types
	config.MountType = stringAsker(`
Would you like to enable netfs for faster filesystem access (y/n)?
(we highly recommend using this option, but this will prompt for password)

(recommended y)
Answer: `, map[string]string{"y": "netfs", "n": "native"})

	config.Save()
	return nil

}

func stringAsker(text string, answers map[string]string) string {
	var answer string

	fmt.Fprint(os.Stdout, text)
	fmt.Scanln(&answer)

	result, ok := answers[answer]
	for !ok {
		fmt.Println("Invalid response, please try again:")
		fmt.Print(text)
		fmt.Scanln(&answer)
		result, ok = answers[answer]
	}
	return result
}

func intAsker(text string, max int) int {
	var answer int

	fmt.Print(text)
	fmt.Scanln(&answer)

	for answer > max {
		fmt.Println("\nInvalid response, please try again:\n")
		fmt.Print(text)
		fmt.Scanln(&answer)
	}
	return answer
}
