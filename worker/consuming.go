package worker

import (
	"fmt"
	"os"
	"time"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"sync"
)

var (
	wg  sync.WaitGroup
	log = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
)

func ConsumeRevisionRun(revisionRun map[int][]string) error {

	clusterConfig, clusterConnection := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))
	log.Info("Connected " + clusterConnection + " the cluster")

	//CREATE AND WATCH PRS
	for i := 0; i < (len(revisionRun)); i++ {

		log.Info("Stage: ", i)

		for j, pr := range revisionRun[i] {

			wg.Add(1)
			log.Info("Concurrent pipelines for this stage:", len(revisionRun[i]))

			fmt.Println(pr)

			renderedPipelineRun := pr
			stage := i
			pipeline := j

			resourceName, _ := sthingsBase.GetRegexSubMatch(renderedPipelineRun, `name: "(.*?)"`)

			go func() {

				defer wg.Done()
				sthingsK8s.CreateDynamicResourcesFromTemplate(clusterConfig, []byte(renderedPipelineRun), tektonNamespace)

				time.Sleep(10 * time.Second)

				log.Info("Stage: ", stage)
				log.Info("Pipeline: ", pipeline)
				log.Info("Verify for: ", resourceName)

				// VerifyPipelineRunStatus(resourceName)

			}()

		}
		wg.Wait()

	}

	fmt.Println("I am here")
	// os.Exit(0)

	return nil
}