package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

//Global variable of Configuration
var (
	dbName              string = "test"
	username            string = "heronalps"
	password            string = "123456"
	dbIP                string = "127.0.0.1"
	dbPort              int    = 3306
	namespace           string = "racelab"
	RunDeployment       string = "image-clf-inf"
	timeQueryDeployment string = "image-clf-inf37"
	// runtimes                   = map[string]bool{"edge": true, "cpu": true, "gpu1": true, "gpu2": true}
	runtimes = map[string]bool{"cpu": true}
	// runtimes               = []string{"gpu1", "gpu2"}
	// NautilusRuntimes       = map[string]bool{"cpu": true, "gpu1": true, "gpu2": true}
	setupImageNums []int = []int{33, 10}
	// currentRuntime      string = "cpu" // cpu / gpu1 / gpu2
	procTimeNumDP        int    = 10
	serviceAccountConfig string = "KUBECONFIG=~/.kube/service-account"
	minikubeConfig       string = "KUBECONFIG=~/.kube/config"
	invokeFile           string = "./scripts/invoke_inf.sh"
	defaultDeployTimes          = map[string]float64{"cpu": 18.0, "gpu1": 46.0, "gpu2": 65.0}
	deploymentTimeNumDP  int    = 1000
	maxWinSize           int    = 100
	minWinSize           int    = 1
	windowSizes                 = map[string]int{"cpu": 1, "gpu1": 1, "gpu2": 1}
	// Workload contains a month worth of randomly generated workload
	Workload      []int  = []int{11, 130, 174, 56, 66, 28, 33, 53, 45, 81, 107, 28, 35, 27, 61, 74, 106, 29, 22, 10, 51, 74, 111, 55, 4, 1, 48, 88, 90, 74, 94, 19, 48, 46, 97, 68, 76, 45, 67, 83, 101, 27, 42, 9, 2, 11, 31, 45, 104, 98, 3, 14, 30, 41, 70, 9, 1, 32, 33, 50, 154, 73, 38, 34, 37, 54, 85, 18, 21, 20, 32, 63, 6, 101, 22, 26, 48, 44, 76, 40, 21, 26, 35, 88, 144, 20, 12, 8, 63, 59, 108, 13, 4, 45, 42, 82, 14, 3, 37, 34, 93, 125, 119, 6, 6, 62, 60, 86, 2, 10, 42, 38, 62, 1, 13, 35, 27, 48, 119, 151, 2, 8, 35, 85, 93, 7, 11, 24, 34, 59, 16, 93, 17, 23, 45, 145, 66, 67, 37, 29, 92, 103, 104, 32, 17, 14, 59, 80, 100, 7, 9, 45, 61, 56, 61, 16, 40, 74, 85, 114, 144, 16, 3, 59, 85, 79, 7, 3, 41, 48, 58, 13, 89, 13, 32, 45, 133, 119, 3, 22, 50, 50, 79, 5, 9, 36, 31, 58, 60, 77, 31, 39, 44, 51, 50, 17, 22, 27, 40, 71, 70, 22, 16, 23, 37, 32, 60, 61, 64, 14, 1, 28, 27, 57, 129, 110, 34, 15, 36, 85, 93, 10, 11, 28, 65, 61, 80, 96, 16, 60, 46, 89, 118, 22, 36, 27, 70, 85, 106, 26, 20, 61, 55, 80, 17, 12, 44, 42, 92, 60, 37, 22, 36, 60, 81, 112, 18, 11, 41, 52, 83, 13, 1, 31, 39, 53, 68, 65, 31, 1, 58, 112, 105, 8, 2, 38, 87, 77, 12, 8, 31, 22, 57, 127, 74, 24, 31, 46, 83, 120, 20, 20, 11, 24, 21, 50, 63, 3, 34, 79, 8, 3, 32, 66, 126, 3, 46, 93, 105, 11, 34, 69, 16, 35, 7, 46, 111, 4, 48, 72, 10, 16, 40, 84, 36, 9, 63, 89, 85, 14, 36, 64, 46, 21, 8, 81, 101, 14, 44, 78, 41, 29, 41, 97, 52, 38, 44, 7, 87, 95, 3, 47, 55, 17, 18, 79, 72, 101, 13, 46, 79, 19, 26, 38, 95, 34, 24, 55, 77, 82, 19, 45, 138, 29, 39, 66, 66, 17, 20, 55, 124, 34, 32, 52, 84, 2, 27, 60, 68, 30, 31, 44, 97, 4, 34, 72, 11, 52, 31, 50, 124, 26, 50, 45, 94, 5, 39, 53, 77, 7, 23, 72, 88, 3, 38, 67, 67, 9, 9, 4, 33, 68, 75, 54, 87, 56, 48, 144, 128}
	presetZipPath string = "/racelab/image_batch_2.zip"

	s3Config = &aws.Config{
		Credentials:      credentials.NewSharedCredentials("", "prp"),
		Endpoint:         aws.String("https://s3.nautilus.optiputer.net"),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	}
)

/*
TimeLog contains total response time and its three components
*/
type TimeLog struct {
	Total      float64
	Transfer   float64
	Deployment float64
	Processing float64
}

/*
DeployResult contains progressed and timeout boolean values
*/
type DeployResult struct {
	Progressed bool
	Timeout    bool
}
