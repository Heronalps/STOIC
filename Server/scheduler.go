/*
Scheduler decides where to run WTB inferencing job given Pi bandwidth, number of images
*/

package server

import "fmt"

/*
GetTime obtains total times of four scenarios
*/
func GetTime() {
	imageNum := ImageCache()
	totalTimes := GetTotalTime(imageNum)
	fmt.Println(totalTimes)
}
