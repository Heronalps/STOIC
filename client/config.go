package client

//Global variable of Configuration
var (
	dbName              string   = "test"
	username            string   = "root"
	password            string   = "123456"
	ip                  string   = "127.0.0.1"
	port                int      = 3306
	namespace           string   = "racelab"
	deployment          string   = "image-clf-inf"
	runtimes            []string = []string{"edge", "cpu", "gpu1", "gpu2"}
	setupImageNums      []int    = []int{33, 10}
	currentRuntime      string   = "edge"
	procTimeNumDP       int      = 10
	deploymentTimeNumDP int      = 10
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
