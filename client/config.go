package client

//Global variable of Configuration
var (
	dbName              string   = "test"
	username            string   = "root"
	password            string   = "123456"
	dbIP                string   = "127.0.0.1"
	dbPort              int      = 3306
	namespace           string   = "racelab"
	deployment          string   = "image-clf-inf"
	timeQueryDeployment string   = "image-clf-inf37"
	runtimes            []string = []string{"edge", "cpu", "gpu1", "gpu2"}
	NautilusRuntimes    map[string]struct{}
	setupImageNums      []int  = []int{33, 10}
	currentRuntime      string = "cpu" // cpu / gpu1 / gpu2
	procTimeNumDP       int    = 10
	deploymentTimeNumDP int    = 10
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
