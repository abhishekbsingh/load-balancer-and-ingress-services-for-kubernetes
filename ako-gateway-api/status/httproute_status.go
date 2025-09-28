/*
 * Copyright © 2025 Broadcom Inc. and/or its subsidiaries. All Rights Reserved.
 * All Rights Reserved.
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*   http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

package status

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"

	akogatewayapilib "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/ako-gateway-api/lib"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/lib"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/status"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/utils"
)

type httproute struct{}

func (o *httproute) Get(key string, name string, namespace string) *gatewayv1.HTTPRoute {

	obj, err := akogatewayapilib.AKOControlConfig().GatewayApiInformers().HTTPRouteInformer.Lister().HTTPRoutes(namespace).Get(name)
	if err != nil {
		utils.AviLog.Warnf("key: %s, msg: unable to get the HTTPRoute object. err: %s", key, err)
		return nil
	}
	utils.AviLog.Debugf("key: %s, msg: Successfully retrieved the HTTPRoute object %s", key, name)
	return obj.DeepCopy()
}

func (o *httproute) GetAll(key string) map[string]*gatewayv1.HTTPRoute {

	objs, err := akogatewayapilib.AKOControlConfig().GatewayApiInformers().HTTPRouteInformer.Lister().List(labels.Everything())
	if err != nil {
		utils.AviLog.Warnf("key: %s, msg: unable to get the HTTPRoute objects. err: %s", key, err)
		return nil
	}

	httpRouteMap := make(map[string]*gatewayv1.HTTPRoute)
	for _, obj := range objs {
		httpRouteMap[obj.Namespace+"/"+obj.Name] = obj.DeepCopy()
	}

	utils.AviLog.Debugf("key: %s, msg: Successfully retrieved the HTTPRoute objects", key)
	return httpRouteMap
}

func (o *httproute) Delete(key string, option status.StatusOptions) {
	nsName := strings.Split(option.Options.ServiceMetadata.HTTPRoute, "/")
	if len(nsName) != 2 {
		utils.AviLog.Warnf("key: %s, msg: invalid HttpRoute name and namespace", key)
		return
	}
	namespace := nsName[0]
	name := nsName[1]
	httpRoute := o.Get(key, name, namespace)
	if httpRoute != nil {
		if option.Options.ServiceMetadata.HTTPRouteRuleName != "" {
			// Remove specific backend entry
			if err := o.patchHTTPRouteAnnotationsDelete(key, httpRoute, option.Options.ServiceMetadata.HTTPRouteRuleName, option.Options.ServiceMetadata.Gateway); err != nil {
				utils.AviLog.Errorf("key: %s, msg: failed to remove HTTPRoute %s/%s annotations for rule %s: %v", key, namespace, name, option.Options.ServiceMetadata.HTTPRouteRuleName, err)
			}
		}
	}
}

func (o *httproute) Update(key string, option status.StatusOptions) {
	nsName := strings.Split(option.Options.ServiceMetadata.HTTPRoute, "/")
	if len(nsName) != 2 {
		utils.AviLog.Warnf("key: %s, msg: invalid HttpRoute name and namespace", key)
		return
	}
	namespace := nsName[0]
	name := nsName[1]
	httpRoute := o.Get(key, name, namespace)
	if httpRoute != nil {
		if option.Options.Status != nil {
			o.Patch(key, httpRoute, option.Options.Status)
		}
		// Update HTTPRoute annotations using patch
		if option.Options.VirtualServiceUUID != "" && option.Options.ServiceMetadata.HTTPRouteRuleName != "" {
			if err := o.patchHTTPRouteAnnotations(key, httpRoute, option.Options.ServiceMetadata.HTTPRouteRuleName, option.Options.ServiceMetadata.Gateway, option.Options.VirtualServiceUUID); err != nil {
				utils.AviLog.Errorf("key: %s, msg: failed to update HTTPRoute annotations: %v", key, err)
			}
		}
	}
}

func (o *httproute) BulkUpdate(key string, options []status.StatusOptions) {
	// TODO: Add this code when we publish the status from the rest layer
}

func (o *httproute) Patch(key string, obj runtime.Object, status *status.Status, retryNum ...int) error {
	retry := 0
	if len(retryNum) > 0 {
		retry = retryNum[0]
		if retry >= 5 {
			utils.AviLog.Errorf("key: %s, msg: Patch retried 5 times, aborting", key)
			akogatewayapilib.AKOControlConfig().EventRecorder().Eventf(obj, corev1.EventTypeWarning, lib.PatchFailed, "Patch of status failed after multiple retries")
			return errors.New("Patch retried 5 times, aborting")
		}
	}

	httpRoute := obj.(*gatewayv1.HTTPRoute)
	if o.isStatusEqual(&httpRoute.Status, status.HTTPRouteStatus) {
		return nil
	}

	patchPayload, _ := json.Marshal(map[string]interface{}{
		"status": status.HTTPRouteStatus,
	})
	_, err := akogatewayapilib.AKOControlConfig().GatewayAPIClientset().GatewayV1().HTTPRoutes(httpRoute.Namespace).Patch(context.TODO(), httpRoute.Name, types.MergePatchType, patchPayload, metav1.PatchOptions{}, "status")
	if err != nil {
		utils.AviLog.Warnf("key: %s, msg: there was an error in updating the HTTPRoute status. err: %+v, retry: %d", key, err, retry)
		updatedObj, err := akogatewayapilib.AKOControlConfig().GatewayApiInformers().HTTPRouteInformer.Lister().HTTPRoutes(httpRoute.Namespace).Get(httpRoute.Name)
		if err != nil {
			utils.AviLog.Warnf("HTTPRoute not found %v", err)
			return err
		}
		return o.Patch(key, updatedObj, status, retry+1)
	}

	utils.AviLog.Infof("key: %s, msg: Successfully updated the HTTPRoute %s/%s status %+v", key, httpRoute.Namespace, httpRoute.Name, utils.Stringify(status))
	return nil
}

func (o *httproute) isStatusEqual(old, new *gatewayv1.HTTPRouteStatus) bool {
	oldStatus, newStatus := old.DeepCopy(), new.DeepCopy()
	currentTime := metav1.Now()
	for i := range oldStatus.Parents {
		for j := range oldStatus.Parents[i].Conditions {
			oldStatus.Parents[i].Conditions[j].LastTransitionTime = currentTime
		}
	}
	for i := range newStatus.Parents {
		for j := range newStatus.Parents[i].Conditions {
			newStatus.Parents[i].Conditions[j].LastTransitionTime = currentTime
		}
	}
	return reflect.DeepEqual(oldStatus, newStatus)
}

func (o *httproute) patchHTTPRouteAnnotations(key string, httpRoute *gatewayv1.HTTPRoute, ruleName, gatewayNSName, virtualServiceUUID string, retryNum ...int) error {
	retry := 0
	if len(retryNum) > 0 {
		retry = retryNum[0]
		if retry >= 3 {
			return errors.New("patchHTTPRouteAnnotations retried 3 times, aborting")
		}
	}

	annotations := httpRoute.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}

	gatewayParts := strings.Split(gatewayNSName, "/")
	if len(gatewayParts) != 2 {
		return fmt.Errorf("invalid gateway namespace name: %s", gatewayNSName)
	}

	// Get existing VS UUID map or create new one
	annotationKey := fmt.Sprintf(lib.HTTPRouteVSAnnotation, gatewayParts[0], gatewayParts[1], ruleName)
	annotations[annotationKey] = virtualServiceUUID

	patchPayload := map[string]interface{}{
		"metadata": map[string]map[string]string{
			"annotations": annotations,
		},
	}

	patchPayloadBytes, err := json.Marshal(patchPayload)
	if err != nil {
		return fmt.Errorf("error marshalling HTTPRoute annotation patch payload: %v", err)
	}

	_, err = akogatewayapilib.AKOControlConfig().GatewayAPIClientset().GatewayV1().HTTPRoutes(httpRoute.Namespace).Patch(context.TODO(), httpRoute.Name, types.MergePatchType, patchPayloadBytes, metav1.PatchOptions{})
	if err != nil {
		utils.AviLog.Warnf("key: %s, msg: error updating HTTPRoute annotations, retry: %d, err: %v", key, retry, err)
		// Fetch updated HTTPRoute and retry
		updatedHTTPRoute, fetchErr := akogatewayapilib.AKOControlConfig().GatewayApiInformers().HTTPRouteInformer.Lister().HTTPRoutes(httpRoute.Namespace).Get(httpRoute.Name)
		if fetchErr != nil {
			return fmt.Errorf("failed to fetch updated HTTPRoute: %v", fetchErr)
		}
		return o.patchHTTPRouteAnnotations(key, updatedHTTPRoute, ruleName, gatewayNSName, virtualServiceUUID, retry+1)
	}

	utils.AviLog.Infof("key: %s, msg: Successfully updated HTTPRoute %s/%s annotations with VS UUID: %s", key, httpRoute.Namespace, httpRoute.Name, virtualServiceUUID)
	return nil
}

func (o *httproute) patchHTTPRouteAnnotationsDelete(key string, httpRoute *gatewayv1.HTTPRoute, ruleName, gatewayNSName string, retryNum ...int) error {
	retry := 0
	if len(retryNum) > 0 {
		retry = retryNum[0]
		if retry >= 3 {
			return errors.New("patchHTTPRouteAnnotationsDelete retried 3 times, aborting")
		}
	}

	annotations := httpRoute.Annotations
	if annotations == nil {
		utils.AviLog.Debugf("key: %s, msg: No annotations found for HTTPRoute %s/%s", key, httpRoute.Namespace, httpRoute.Name)
		return nil
	}

	gatewayParts := strings.Split(gatewayNSName, "/")
	if len(gatewayParts) != 2 {
		return fmt.Errorf("invalid gateway namespace name: %s", gatewayNSName)
	}

	annotationKey := fmt.Sprintf(lib.HTTPRouteVSAnnotation, gatewayParts[0], gatewayParts[1], ruleName)

	_, exists := annotations[annotationKey]
	if !exists {
		utils.AviLog.Debugf("key: %s, msg: %s annotation not found for HTTPRoute %s/%s", key, annotationKey, httpRoute.Namespace, httpRoute.Name)
		return nil
	}
	delete(annotations, annotationKey)

	patchPayload := map[string]interface{}{
		"metadata": map[string]map[string]string{
			"annotations": annotations,
		},
	}

	patchPayloadBytes, err := json.Marshal(patchPayload)
	if err != nil {
		return fmt.Errorf("error marshalling HTTPRoute annotation delete patch payload: %v", err)
	}

	_, err = akogatewayapilib.AKOControlConfig().GatewayAPIClientset().GatewayV1().HTTPRoutes(httpRoute.Namespace).Patch(context.TODO(), httpRoute.Name, types.MergePatchType, patchPayloadBytes, metav1.PatchOptions{})
	if err != nil {
		utils.AviLog.Warnf("key: %s, msg: error removing HTTPRoute annotations, retry: %d, err: %v", key, retry, err)
		// Fetch updated HTTPRoute and retry
		updatedHTTPRoute, fetchErr := akogatewayapilib.AKOControlConfig().GatewayApiInformers().HTTPRouteInformer.Lister().HTTPRoutes(httpRoute.Namespace).Get(httpRoute.Name)
		if fetchErr != nil {
			return fmt.Errorf("failed to fetch updated HTTPRoute: %v", fetchErr)
		}
		return o.patchHTTPRouteAnnotationsDelete(key, updatedHTTPRoute, ruleName, gatewayNSName, retry+1)
	}

	utils.AviLog.Infof("key: %s, msg: Successfully removed HTTPRoute %s/%s VS UUID annotation", key, httpRoute.Namespace, httpRoute.Name)
	return nil
}
