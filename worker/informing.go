package worker

import "os"

var (
	logfilePath     = "yaw-informer.log"
	prApiName       = "tekton.dev"
	prApiVersion    = "v1beta1"
	prApiResource   = "pipelineruns"
	tektonNamespace = os.Getenv("TEKTON_NAMESPACE")
)
