// The server scheduler program running on Sedgwick Euca.
// It solicits the bandwidth of camera (Pi Zero) and decides to run WTB locally
// or transfer images to Mayhem cloud and run WTb at Nautilus.

package server

import (
	"bufio"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

/*
GetBandWidth solicits bandwidth of Pi Zero at Sedgwick Reserve.
*/
func GetBandWidth() float64 {
	url := "http://169.231.235.221/sedgtomayhem.txt"
	lines := "1"
	fmt.Printf("Tailing http endpoint : %s\n", url)

	cmd := "curl " + url + " | tail -" + lines
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		println(err.Error())
		return 0
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var lastLine string
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		lastLine = scanner.Text()
	}
	fields := strings.Fields(lastLine)
	bandWidth, err := strconv.ParseFloat(fields[0], 64)
	return bandWidth
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
