/*
Scheduler decides where to run WTB inferencing job given Pi bandwidth, number of images
*/

package client

import (
	"fmt"
	exec "os/exec"
	"sort"
	"strconv"
)

/*
Schedule is the entry point of Scheduler.
Parameter: runtime is the preset runtime. If it's empty string, SelectRunTime will select one.
Return: runtime for appending processing time to corresponding table
*/
func Schedule(runtime string, imageNum int, app string, version string) []byte {
	fmt.Printf("The bandwidth is %f megabits \n", GetBandWidth())
	fmt.Printf("The batch of %d images needs %f seconds to transfer\n", imageNum, GetTransferTime(imageNum))

	var (
		elapsed float64
		output  []byte
	)
	if runtime == "" {
		runtime = SelectRunTime(imageNum, app, version)
		// Update current runtime to accurately estimate deployment time
		currentRuntime = runtime
	}
	output, elapsed = Request(runtime, imageNum, app, version)
	if elapsed != 0.0 {
		elapsed += GetAdditionTime(runtime, imageNum)
		AppendRecordProcessing(dbName, runtime, imageNum, elapsed, app, version)
	}
	return output
}

/*
Request is a wrap function both for executing jobs and setting up processing time table for regression
*/
func Request(runtime string, imageNum int, app string, version string) ([]byte, float64) {
	var (
		output  []byte
		elapsed float64
	)
	switch runtime {
	case "edge":
		fmt.Println("Running on edge...")
		output, elapsed = RunOnEdge(imageNum, app, version)
	default:
		fmt.Println("Running on Nautilus...")
		output, elapsed = RunOnNautilus(runtime, imageNum, app, version)
	}
	return output, elapsed
}

/*
SelectRunTime select the runtime among four scenarios
*/
func SelectRunTime(imageNum int, app string, version string) string {
	totalTimes := GetTotalTime(imageNum, app, version)
	fmt.Println(totalTimes)
	// Sort the totalTimes map by key
	keys := make([]float64, 0, len(totalTimes))
	for k := range totalTimes {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	fmt.Printf("The task is scheduled at %s for %f seconds\n", totalTimes[keys[0]], keys[0])
	return totalTimes[keys[0]]
}

/*
RunOnEdge runs the task on mini edge cloud with AVX support
*/
func RunOnEdge(imageNum int, app string, version string) ([]byte, float64) {
	var (
		output []byte
		err    error
		cmd    *exec.Cmd
	)
	repoPATH := HomeDir() + "/GPU_Serverless"

	// Run WTB image classification task
	FILE := "./kubeless/image_clf/inference/local_version/image_clf_inf.py "
	cmdRun := "source venv/bin/activate && python " + FILE + strconv.Itoa(int(imageNum))
	cmd = exec.Command("bash", "-c", cmdRun)
	cmd.Dir = repoPATH
	fmt.Printf("Start running task %s version %s on %d images \n", app, version, imageNum)
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error running task. msg: %s \n", err.Error())
		return output, 0
	}
	fmt.Printf("Output of task %s\n", string(output))
	return output, ParseElapsed(output)
}
