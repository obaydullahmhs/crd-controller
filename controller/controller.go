package controller

import (
	klientset "github.com/obaydullahmhs/crd-controller/pkg/client/clientset/versioned"
	informer "github.com/obaydullahmhs/crd-controller/pkg/client/informers/externalversions/aadee.apps/v1alpha1"
	lister "github.com/obaydullahmhs/crd-controller/pkg/client/listers/aadee.apps/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
)

type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// myclientset is a clientset for our own API group
	myclientset klientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	aadeeLister       lister.AadeeLister
	aadeeSynced       cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workQueue workqueue.RateLimitingInterface
}

// NewController returns a new sample controller
func NewController(kubeclientset kubernetes.Interface,
	myclientset klientset.Interface,
	deploymentInformer appsinformer.DeploymentInformer,
	aadeeInformer informer.AadeeInformer) *Controller {

	ctrl := &Controller{
		kubeclientset:     kubeclientset,
		myclientset:       myclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		aadeeLister:       aadeeInformer.Lister(),
		aadeeSynced:       aadeeInformer.Informer().HasSynced,
		workQueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Aadeez"),
	}

	log.Println("Setting up event handlers")

	// Set up an event handler for when Aadee resources change
	aadeeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: ctrl.enqueueAadeez,
		UpdateFunc: func(oldObj, newObj interface{}) {
			ctrl.enqueueAadeez(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			ctrl.enqueueAadeez(obj)
		},
	})

	// whatif deployment resources changes??

	return ctrl
}

func (c *Controller) enqueueAadeez(obj interface{}) {
	log.Println("Enqueueing Aadeez ...")
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
	}
	c.workQueue.AddRateLimited(key)
}
