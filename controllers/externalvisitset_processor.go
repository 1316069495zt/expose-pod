package controllers

import (
	"context"
	appsv1alpha1 "external/api/v1alpha1"
	sidecarcontrol "external/sidecarcontrol"
	util "external/util"
	"fmt"
	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Processor struct {
	Client   client.Client
	recorder record.EventRecorder
}

func NewSidecarSetProcessor(cli client.Client, rec record.EventRecorder) *Processor {
	return &Processor{
		Client:   cli,
		recorder: rec,
	}
}

func (p *Processor) UpdateExternalvisitSet(externalvisitSet *appsv1alpha1.ExternalvisitSet) (ctrl.Result, error) {

	control := New(externalvisitSet)
	if !control.IsActiveExternalvisitSet() {
		return reconcile.Result{}, nil
	}
	pods, err := p.getMatchingPods(externalvisitSet)

	if err != nil {
		klog.Errorf("externalvisitSet get matching pods error, err: %v, name: %s", err, externalvisitSet.Name)
		return reconcile.Result{}, err
	}

	for i, pod := range pods {
		log.Info("过滤出的pod",i, pod.Name )
	}

	if pods == nil {
		return reconcile.Result{}, nil
	}

	for _, pod := range pods {
		svc := NewService(pod)

		if err := p.Client.Create(context.TODO(), svc); err != nil {

			log.Info("创建svc失败",err.Error() )

			//todo 不操作
			//return ctrl.Result{}, err
		}
	}

	return reconcile.Result{}, nil

}

func NewService(pod *corev1.Pod) *corev1.Service {

	log.Info("给 pod:",pod.Name,"构造svc" )

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(pod, schema.GroupVersionKind{
					Group:   appsv1alpha1.GroupVersion.Group,
					Version: appsv1alpha1.GroupVersion.Version,
					Kind:    pod.Kind,
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			//Ports: pod.Spec.ContainersPorts,
			Ports: []corev1.ServicePort{{
				Protocol: "TCP",
				Port:     80,
				TargetPort: intstr.FromInt(80),
				NodePort: 32333,
			}},
			Selector: map[string]string{
				"app": pod.ObjectMeta.Labels["app"],
			},
		},
	}
}

func (p *Processor) getMatchingPods(s *appsv1alpha1.ExternalvisitSet) ([]*corev1.Pod, error) {
	// get more faster selector
	selector, err := util.GetFastLabelSelector(s.Spec.Selector)
	if err != nil {
		return nil, err
	}
	log.Info("selector:  ", selector)

	if selector == nil {
		return nil, nil
	}

	// If externalvisitSet.Spec.Namespace is empty, then select in cluster
	//scopedNamespaces := []string{s.Spec.Namespace}

	log.Info("开始执行 getMatchingPods ")
	//log.Info(scopedNamespaces)
	selectedPods, err := p.getSelectedPods(nil, selector)
	if err != nil {
		return nil, err
	}
	for i, pod := range selectedPods {
		log.Info("查找出的pod",i, pod.Name )
	}

	// filter out pods that don't require updated, include the following:
	// 1. Deletion pod
	// 2. ignore namespace: "kube-system", "kube-public"
	var filteredPods []*corev1.Pod
	for _, pod := range selectedPods {
		if sidecarcontrol.IsActivePod(pod) {
			filteredPods = append(filteredPods, pod)
		}
	}
	return filteredPods, nil
}

func (p *Processor) getSelectedPods(namespaces []string, selector labels.Selector) (relatedPods []*corev1.Pod, err error) {
	listOpts := &client.ListOptions{LabelSelector: selector}
	//for _, ns := range namespaces {
	log.Info("listOpts : ~~~~~~~~~~~~~~~~~~~~~~~ ", listOpts)

	allPods := &corev1.PodList{}
	//allPods = new(corev1.PodList)
	log.Info("allPods : ~~~~~~~~~~~~~~~~~~~~~~~ ", allPods)

	//listOpts.Namespace = "default"
	if listErr := p.Client.List(context.TODO(), allPods, listOpts); listErr != nil {
		err = fmt.Errorf("externalvisitSet list pods by ns error, , err:%v", listErr)
		return
	}

	for i := range allPods.Items {
		relatedPods = append(relatedPods, &allPods.Items[i])
	}
	//}
	return
}
