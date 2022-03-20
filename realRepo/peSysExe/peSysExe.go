package peSysExe

import (
	"log"
	"os/exec"
)

func PeSysCmd(myCommand string) string {
	cmd := exec.Command("/bin/sh", "-c", myCommand)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	myOutput, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(myOutput)
}

func PeSysCmdWait(myCommand string) string {
	myexecute := exec.Command("/bin/sh", "-c", myCommand+" &")
	myOutput, err := myexecute.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(myOutput)
}
