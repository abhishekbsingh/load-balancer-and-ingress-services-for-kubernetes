module github.com/vmware/load-balancer-and-ingress-services-for-kubernetes

go 1.24.0

toolchain go1.24.3

require (
	github.com/Masterminds/semver v1.5.0
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/fatih/structs v1.1.0
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-logr/logr v1.4.2
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/jinzhu/copier v0.3.5
	github.com/jupp0r/go-priority-queue v0.0.0-20160601094913-ab1073853bde
	github.com/onsi/gomega v1.35.1
	github.com/openshift/api v0.0.0-20201019163320-c6a5ec25f267
	github.com/openshift/client-go v0.0.0-20201020082437-7737f16e53fc
	github.com/prometheus/client_golang v1.17.0
	github.com/vmware-tanzu/service-apis v0.0.0-20200901171416-461d35e58618
	github.com/vmware/alb-sdk v0.0.0-20250407104238-8f88209a8f1a
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.37.0
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
	google.golang.org/protobuf v1.36.6
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	istio.io/client-go v1.25.2
	k8s.io/api v0.33.1
	k8s.io/apimachinery v0.33.1
	k8s.io/client-go v0.33.1
	k8s.io/klog/v2 v2.130.1
	k8s.io/utils v0.0.0-20250321185631-1f6e0b77f77e
	sigs.k8s.io/gateway-api v1.3.0
	sigs.k8s.io/service-apis v0.1.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/term v0.31.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	istio.io/api v1.25.2-0.20250410212420-84c271001f68 // indirect
	k8s.io/kube-openapi v0.0.0-20250318190949-c8a335a9a2ff // indirect
	sigs.k8s.io/json v0.0.0-20241010143419-9aa6b5e7a4b3 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.7.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace (
	github.com/davecgh/go-spew => github.com/davecgh/go-spew v1.1.1
	github.com/golang/glog => github.com/golang/glog v1.2.4
	github.com/google/gofuzz => github.com/google/gofuzz v1.2.0
	github.com/onsi/gomega => github.com/onsi/gomega v1.14.0
	sigs.k8s.io/service-apis => sigs.k8s.io/service-apis v0.1.0
)
