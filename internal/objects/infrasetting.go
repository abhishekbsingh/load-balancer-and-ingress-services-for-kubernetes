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

package objects

import (
	"sync"
)

var infral7lister *AviInfraSettingL7Lister
var infraonce sync.Once

func InfraSettingL7Lister() *AviInfraSettingL7Lister {
	infraonce.Do(func() {
		infral7lister = &AviInfraSettingL7Lister{
			IngRouteInfraSettingStore:  NewObjectMapStore(),
			InfraSettingShardSizeStore: NewObjectMapStore(),
			NSScopedInfraSettingStore:  NewObjectMapStore(),
		}
	})
	return infral7lister
}

type AviInfraSettingL7Lister struct {
	InfraSettingIngRouteLock sync.RWMutex

	// namespaced ingress/route -> infrasetting
	IngRouteInfraSettingStore *ObjectMapStore

	// infrasetting -> shardSize
	InfraSettingShardSizeStore *ObjectMapStore

	// infrasettig -> namespaces
	NSScopedInfraSettingStore *ObjectMapStore
}

func (v *AviInfraSettingL7Lister) GetIngRouteToInfraSetting(ingrouteNsName string) (bool, string) {
	found, infraSettingName := v.IngRouteInfraSettingStore.Get(ingrouteNsName)
	if !found {
		return false, ""
	}
	return true, infraSettingName.(string)
}

func (v *AviInfraSettingL7Lister) UpdateIngRouteInfraSettingMappings(ingrouteNsName, infraSettingName, shardSize string) {
	v.InfraSettingIngRouteLock.Lock()
	defer v.InfraSettingIngRouteLock.Unlock()
	v.IngRouteInfraSettingStore.AddOrUpdate(ingrouteNsName, infraSettingName)
	v.InfraSettingShardSizeStore.AddOrUpdate(infraSettingName, shardSize)
}

func (v *AviInfraSettingL7Lister) RemoveIngRouteInfraSettingMappings(ingrouteNsName string) bool {
	v.InfraSettingIngRouteLock.Lock()
	defer v.InfraSettingIngRouteLock.Unlock()
	mappingDeleted := false
	if found, infraSettingName := v.GetIngRouteToInfraSetting(ingrouteNsName); found {
		// first delete the ingress-infrasetting mapping entry
		mappingDeleted = v.IngRouteInfraSettingStore.Delete(ingrouteNsName)
		// delete infrasetting only if it is not mapped to any other ingress
		if !v.IngRouteInfraSettingStore.IsInfraSettingMapped(infraSettingName) {
			v.InfraSettingShardSizeStore.Delete(infraSettingName)
		}
	}
	return mappingDeleted
}

func (v *AviInfraSettingL7Lister) GetInfraSettingToShardSize(infraSettingName string) (bool, string) {
	found, shardSize := v.InfraSettingShardSizeStore.Get(infraSettingName)
	if !found {
		return false, ""
	}
	return true, shardSize.(string)
}

func (v *AviInfraSettingL7Lister) UpdateInfraSettingToNamespaceMapping(infraSetting string, namespaces []interface{}) {
	v.NSScopedInfraSettingStore.AddOrUpdate(infraSetting, namespaces)
}

func (v *AviInfraSettingL7Lister) GetInfraSettingScopedNamespaces(infraSetting string) []interface{} {
	found, namespaces := v.NSScopedInfraSettingStore.Get(infraSetting)
	if !found {
		return []interface{}{}
	}
	return namespaces.([]interface{})
}

func (v *AviInfraSettingL7Lister) DeleteInfraSettingToNamespaceMapping(infraSetting string) {
	v.NSScopedInfraSettingStore.Delete(infraSetting)
}
