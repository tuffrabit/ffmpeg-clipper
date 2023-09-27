package common

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func RunSystemCommand(cmd *exec.Cmd) (string, error) {
	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	cmdString := cmd.String()

	log.Printf("Running %v\n", cmdString)

	err := cmd.Run()
	outString := cmdOut.String()
	errString := cmdErr.String()
	if err != nil || errString != "" {
		if outString != "" {
			log.Printf("%v stdout: %v", cmdString, outString)
		}

		if errString != "" {
			log.Printf("%v stderr: %v", cmdString, errString)

			if err == nil {
				err = errors.New(errString)
			}
		}

		log.Println(err)

		return errString, err
	} else {
		if outString != "" {
			log.Printf("%v stdout: %v", cmdString, outString)
		}

		return outString, nil
	}
}

func GetRandomString() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X", b)
}
