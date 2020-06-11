/*
 * Copyright 2019-2020 VMware, Inc.
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
package retry

import (
	"ako/pkg/nodes"

	"github.com/avinetworks/container-lib/utils"
)

func DequeueSlowRetry(vsKey string) {
	// Retrieve the Key and note the time.
	utils.AviLog.Infof("Retrieved the key for slow retry : %s", vsKey)
	sharedQueue := utils.SharedWorkQueue().GetQueueByName(utils.GraphLayer)
	modelName := utils.ADMIN_NS + "/" + vsKey
	nodes.PublishKeyToRestLayer(modelName, "retry", sharedQueue)
}

func DequeueFastRetry(vsKey string) {
	utils.AviLog.Infof("Retrieved the key for fast retry: %s", vsKey)
	sharedQueue := utils.SharedWorkQueue().GetQueueByName(utils.GraphLayer)
	modelName := utils.ADMIN_NS + "/" + vsKey
	nodes.PublishKeyToRestLayer(modelName, "retry", sharedQueue)

}