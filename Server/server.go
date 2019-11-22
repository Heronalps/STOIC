// The server scheduler program running on Sedgwick Euca.
// It solicits the bandwidth of camera (Pi Zero) and decides to run WTB locally
// or transfer images to Mayhem cloud and run WTb at Nautilus.

package server

import (
	"fmt"
)

/*
Cronjob runs background cron job.
*/
func Cronjob() {
	fmt.Println("This is cron job!")
}
