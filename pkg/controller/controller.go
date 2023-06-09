package controller

import (
	"fmt"

	controllerv1alpha1 "github.com/obaydullahmhs/crd-controller/pkg/apis/aadee.apps/v1alpha1"
	klientset "github.com/obaydullahmhs/crd-controller/pkg/client/clientset/versioned"
	informer "github.com/obaydullahmhs/crd-controller/pkg/client/informers/externalversions/aadee.apps/v1alpha1"
	lister "github.com/obaydullahmhs/crd-controller/pkg/client/listers/aadee.apps/v1alpha1"
	"golang.org/x/net/context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
	"time"
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
		workQueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Aadees"),
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

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shut down the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workQueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Println("Starting Controller")
	// Wait for the caches to be synced before starting workers
	log.Println("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.aadeeSynced); !ok {
		return fmt.Errorf("failed to wait for cache to sync")
	}

	log.Println("Starting workers")
	// Launch two workers to process Aadee resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	log.Println("Worker Started")
	<-stopCh
	log.Println("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the workqueue.
func (c *Controller) runWorker() {
	for c.ProcessNextItem() {
	}
}

func (c *Controller) ProcessNextItem() bool {
	obj, shutdown := c.workQueue.Get()

	if shutdown {
		return false
	}
	// We wrap this block in a func, so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off period.
		defer c.workQueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// syncHandler is business logic
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workQueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item, so it does not
		// get queued again until another change happens.
		c.workQueue.Forget(obj)
		log.Printf("successfully synced '%s'\n", key)
		return nil
	}(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return true
	}
	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Aadee resource
// with the current status of the resource.
// implement the business logic here.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Aadee resource with this namespace/name
	aadee, err := c.aadeeLister.Aadees(namespace).Get(name)
	if err != nil {
		// The Aadee resource may no longer exist, in which case we stop processing.
		if errors.IsNotFound(err) {
			// We choose to absorb the error here as the worker would requeue the
			// resource otherwise. Instead, the next time the resource is updated
			// the resource will be queued again.
			fmt.Printf("Aadee '%s' in work queue no longer exists", key)
			return nil
		}
		return err
	}

	deploymentName := aadee.Name + "-" + aadee.Spec.Name
	if aadee.Spec.Name == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		deploymentName = aadee.Name + "-missingname"
		//utilruntime.HandleError(fmt.Errorf("%s : deployment name must be specified", key))
		//return nil
	}

	// Get the deployment with the name specified in Aadee.spec
	deployment, err := c.deploymentsLister.Deployments(namespace).Get(deploymentName + "-depl")
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(aadee.Namespace).Create(context.TODO(), newDeployment(aadee, deploymentName), metav1.CreateOptions{})

	}
	// If an error occurs during Get/Create, we'll requeue the item, so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If this number of the replicas on the Aadee resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if aadee.Spec.Replicas != nil && *aadee.Spec.Replicas != *deployment.Spec.Replicas {
		log.Printf("Aadee %s replicas: %d, deployment replicas: %d\n", name, *aadee.Spec.Replicas, *deployment.Spec.Replicas)

		deployment, err = c.kubeclientset.AppsV1().Deployments(namespace).Update(context.TODO(), newDeployment(aadee, deploymentName), metav1.UpdateOptions{})
		// If an error occurs during Update, we'll requeue the item, so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}
	}
	// Finally, we update the status block of the Aadee resource to reflect the
	// current state of the world
	err = c.updateAadeeStatus(aadee, deployment)
	if err != nil {
		return err
	}

	serviceName := aadee.Name + "-" + aadee.Spec.Name
	if aadee.Spec.Name == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		serviceName = aadee.Name + "-missingname"
		//utilruntime.HandleError(fmt.Errorf("%s : deployment name must be specified", key))
		//return nil
	}
	// Check if service already exists or not
	service, err := c.kubeclientset.CoreV1().Services(aadee.Namespace).Get(context.TODO(), serviceName+"-svc", metav1.GetOptions{})

	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.CoreV1().Services(aadee.Namespace).Create(context.TODO(), newService(aadee, serviceName), metav1.CreateOptions{})
		if err != nil {
			//log.Println(err)
			return err
		}
		log.Printf("\nservice %s created .....\n", service.Name)
	} else if err != nil {
		log.Println(err)
		return err
	}

	_, err = c.kubeclientset.CoreV1().Services(aadee.Namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *Controller) updateAadeeStatus(aadee *controllerv1alpha1.Aadee, deployment *appsv1.Deployment) error {

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	aadeeCopy := aadee.DeepCopy()
	aadeeCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Foo resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.myclientset.AadeeV1alpha1().Aadees(aadee.Namespace).Update(context.TODO(), aadeeCopy, metav1.UpdateOptions{})

	return err

}

// newDeployment creates a new Deployment for a Aadee resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Aadee resource that 'owns' it.
func newDeployment(aadee *controllerv1alpha1.Aadee, deploymentName string) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName + "-depl",
			Namespace: aadee.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "aadee.apps/v1alpha1",
					Kind:       "Aadee",
					Name:       aadee.Name,
					UID:        aadee.UID,
					Controller: func() *bool {
						var ok = true
						return &ok
					}(),
				},
			},
		},
		Spec: appsv1.DeploymentSpec{

			Replicas: aadee.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  aadee.Name,
					"kind": "Aadee",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  aadee.Name,
						"kind": "Aadee",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  aadee.Name,
							Image: aadee.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: aadee.Spec.Container.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func newService(aadee *controllerv1alpha1.Aadee, serviceName string) *corev1.Service {
	labels := map[string]string{
		"app":  aadee.Name,
		"kind": "Aadee",
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName + "-svc",
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(aadee, controllerv1alpha1.SchemeGroupVersion.WithKind("Aadee")),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Port:       aadee.Spec.Container.Port,
					TargetPort: intstr.FromInt(int(aadee.Spec.Container.Port)),
					Protocol:   "TCP",
				},
			},
		},
	}

}
