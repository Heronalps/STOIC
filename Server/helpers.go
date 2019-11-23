/*
This module contains all helper functions.
*/

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

/*
Extrapolate function inferences runtime by coefficient and intercept.
*/
func Extrapolate(mode string, x float64) float64 {
	var coef float64
	var intercept float64
	switch mode {
	case "euca":
		coef = 2.39549861
		intercept = 13.600537473199736
	case "cpu":
		coef = 1.35328625
		intercept = 29.910393951759985
	case "gpu1":
		coef = 0.34666546
		intercept = 92.16208453231344
	case "gpu2":
		coef = 0.21877092
		intercept = 143.26647033799347
	}

	return x*coef + intercept
}
