/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"

	sthingsK8s "github.com/stuttgart-things/sthingsK8s"
)

var (
	logfilePath     = "yw.log"
	prApiName       = "tekton.dev"
	prApiVersion    = "v1beta1"
	prApiResource   = "pipelineruns"
	tektonNamespace = os.Getenv("TEKTON_NAMESPACE")
)

func VerifyPipelineRunStatus(name string) {

	clusterConfig, _ := sthingsK8s.GetKubeConfig(os.Getenv("KUBECONFIG"))

	log.Info("use elasticsearch: " + sendToElastic)

	clusterClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalln(err)
	}

	resource := schema.GroupVersionResource{Group: prApiName, Version: prApiVersion, Resource: prApiResource}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient, time.Minute, tektonNamespace, nil)
	informer := factory.ForResource(resource).Informer()

	mux := &sync.RWMutex{}
	synced := false

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			mux.RLock()

			defer mux.RUnlock()

			if !synced {
				return
			}

			pr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldObj)
			fmt.Println(err)

			var prList *v1beta1.PipelineRun

			err = runtime.DefaultUnstructuredConverter.FromUnstructured(pr, &prList)
			if err != nil {
				log.Fatal(err)
			}

			if prList.Name == name {

				log.Info("found pipelineRun", prList.Name)

				status := fmt.Sprint(prList.Status)

				prStatus := GetPipelineRunStatus(status)

				if prStatus["status"] != "Unknown" {

					log.Info("Status is not unknown. end watching the resource..", prList.Name)

					syscall.Kill(syscall.Getpid(), syscall.SIGINT)

				}

				fmt.Println("STATUS IS UNKNOW", prList.Name)
			}

		},
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go informer.Run(ctx.Done())

	isSynced := cache.WaitForCacheSync(ctx.Done(), informer.HasSynced)
	mux.Lock()
	synced = isSynced
	mux.Unlock()

	if !isSynced {
		log.Fatal("failed to sync")
	}

	<-ctx.Done()

}
