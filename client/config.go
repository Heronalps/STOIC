package client

//Global variable of Configuration
var (
	dbName              string = "test"
	username            string = "root"
	password            string = "Stoic!@#$%^123456"
	dbIP                string = "127.0.0.1"
	dbPort              int    = 3306
	namespace           string = "racelab"
	RunDeployment       string = "image-clf-inf"
	timeQueryDeployment string = "image-clf-inf37"
	runtimes                   = map[string]bool{"edge": true, "cpu": true, "gpu1": true, "gpu2": true}
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
