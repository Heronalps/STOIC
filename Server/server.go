// The server scheduler program running on Sedgwick Euca.
// It solicits the bandwidth of camera (Pi Zero) and decides to run WTB locally
// or transfer images to Mayhem cloud and run WTb at Nautilus.

package server

import (
	"fmt"
	"os/exec"
	"os/user"
)

/*
TailHTTP runs background cron job.
*/
func TailHTTP() {
	url := "http://169.231.235.221/sedgtomayhem.txt"
	htailPath := HomeDir() + "/Downloads/Github/htail/htail.py"
	command := "python2"
	fmt.Printf("Tailing http endpoint : %s\n", url)

	cmd := exec.Command(command, htailPath, url)
	out, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(string(out))
}

/*
HomeDir gets home directory of current user
*/
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		println(err.Error())
		return ""
	}
	return usr.HomeDir
}
