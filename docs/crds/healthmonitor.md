### HealthMonitor

The HealthMonitor CRD is managed by the AKO CRD Operator and allows users to define custom health monitoring configurations for backend services. This CRD enables fine-grained control over health check parameters, supporting various health check types including TCP, PING, and HTTP monitors.

A sample HealthMonitor CRD looks like this:

```yaml
apiVersion: ako.vmware.com/v1alpha1
kind: HealthMonitor
metadata:
  name: my-health-monitor
  namespace: default
spec:
  type: HEALTH_MONITOR_HTTP
  sendInterval: 10
  receiveTimeout: 4
  successfulChecks: 2
  failedChecks: 2
  monitorPort: 80
  isFederated: false
  httpMonitor:
    httpRequest: "GET /health HTTP/1.1"
    httpResponseCode:
      - HTTP_2XX
    httpResponse: "OK"
    exactHttpRequest: false
    maintenanceCode:
      - 503
    maintenanceResponse: "Maintenance"
    responseSize: 2048
```

#### Health Monitor Types

The HealthMonitor CRD supports three types of health monitors, specified using the `type` field:

- `HEALTH_MONITOR_TCP`: TCP-based health check
- `HEALTH_MONITOR_PING`: ICMP ping-based health check  
- `HEALTH_MONITOR_HTTP`: HTTP-based health check

```yaml
  type: HEALTH_MONITOR_HTTP
```

**Note**: The `type` field is required and immutable once the HealthMonitor is created.

#### Basic Health Check Parameters

The following parameters are common across all health monitor types:

##### Send Interval

The `sendInterval` defines the frequency, in seconds, that health checks are sent to backend servers.

```yaml
  sendInterval: 10
```

- Valid range: 1-3600 seconds
- Default: 10 seconds

##### Receive Timeout

The `receiveTimeout` is the timeout for receiving a health check response, in seconds.

```yaml
  receiveTimeout: 4
```

- Valid range: 1-2400 seconds
- Default: 4 seconds

##### Successful Checks

The `successfulChecks` is the number of consecutive successful health checks required before marking a server as UP.

```yaml
  successfulChecks: 2
```

- Valid range: 1-50
- Default: 2

##### Failed Checks

The `failedChecks` is the number of consecutive failed health checks required before marking a server as DOWN.

```yaml
  failedChecks: 2
```

- Valid range: 1-50
- Default: 2

##### Monitor Port

The `monitorPort` specifies the port to use for the health check. If not specified, the service port is used.

```yaml
  monitorPort: 8080
```

- Valid range: 0-65535
- **Note**: This field is immutable once the HealthMonitor is created.

#### TCP Health Monitor Configuration

For TCP-based health monitors (`HEALTH_MONITOR_TCP`), additional configuration can be specified under the `tcpMonitor` section:

```yaml
  type: HEALTH_MONITOR_TCP
  tcpMonitor:
    tcpRequest: "GET /health HTTP/1.1"
    tcpResponse: "OK"
    maintenanceResponse: "Maintenance"
    tcpHalfOpen: false
```

##### TCP Request

The `tcpRequest` is the data to send as part of the TCP health check.

- Maximum length: 1024 characters
- Optional

##### TCP Response

The `tcpResponse` is the expected response from the server to consider the health check successful.

- Maximum length: 512 characters
- Optional

##### Maintenance Response

The `maintenanceResponse` is the response that indicates the server is in maintenance mode.

- Maximum length: 512 characters
- Optional

##### TCP Half Open

The `tcpHalfOpen` flag determines if the TCP monitor should use TCP half-open (SYN-only) checks.

```yaml
  tcpHalfOpen: true
```

**Note**: When `tcpHalfOpen` is set to true, the `tcpRequest`, `tcpResponse`, and `maintenanceResponse` fields must not be set.

#### HTTP Health Monitor Configuration

For HTTP-based health monitors, configuration is specified under the `httpMonitor` section:

```yaml
  type: HEALTH_MONITOR_HTTP
  httpMonitor:
    httpRequest: "GET /health HTTP/1.1"
    httpResponseCode:
      - HTTP_2XX
      - HTTP_3XX
    httpResponse: "healthy"
    exactHttpRequest: false
    maintenanceCode:
      - 503
      - 500
    maintenanceResponse: "Maintenance"
    authType: AUTH_BASIC
    httpRequestBody: '{"check": "health"}'
    responseSize: 4096
```

##### HTTP Request

The `httpRequest` specifies the HTTP request to send for the health check.

```yaml
  httpRequest: "GET /health HTTP/1.1"
```

- Maximum length: 1024 characters
- Default: "GET / HTTP/1.0"

##### HTTP Response Code

The `httpResponseCode` specifies the list of HTTP response code ranges that indicate a successful health check.

```yaml
  httpResponseCode:
    - HTTP_2XX
    - HTTP_3XX
```

Valid values are:
- `HTTP_ANY`: Any HTTP response code
- `HTTP_1XX`: 1xx response codes
- `HTTP_2XX`: 2xx response codes
- `HTTP_3XX`: 3xx response codes
- `HTTP_4XX`: 4xx response codes
- `HTTP_5XX`: 5xx response codes

**Note**: At least one response code must be specified.

##### HTTP Response

The `httpResponse` is a keyword to match in the response body to consider the health check successful.

```yaml
  httpResponse: "OK"
```

- Maximum length: 512 characters
- Optional

##### Exact HTTP Request

The `exactHttpRequest` flag determines if the entire HTTP request should match exactly as specified.

```yaml
  exactHttpRequest: true
```

- Default: false
- **Note**: If `authType` is set, `exactHttpRequest` must be set to false.

##### Maintenance Code

The `maintenanceCode` specifies HTTP response codes that indicate the server is in maintenance mode.

```yaml
  maintenanceCode:
    - 503
    - 500
```

- Valid range for each code: 101-599
- Maximum items: 4
- Optional

##### Maintenance Response

The `maintenanceResponse` specifies body content to match that indicates the server is in maintenance mode.

```yaml
  maintenanceResponse: "Under Maintenance"
```

- Maximum length: 512 characters
- Optional

##### Authentication Type

The `authType` specifies the authentication method to use for HTTP health checks.

```yaml
  authType: AUTH_BASIC
```

Valid values:
- `AUTH_BASIC`: HTTP Basic Authentication
- `AUTH_NTLM`: NTLM Authentication

**Note**: When `authType` is set, the `authentication.secretRef` field must also be configured.

##### HTTP Request Body

The `httpRequestBody` specifies the request body to send with the HTTP health check.

```yaml
  httpRequestBody: '{"check": "health"}'
```

- Optional

##### Response Size

The `responseSize` specifies the expected maximum size of the HTTP response in bytes.

```yaml
  responseSize: 4096
```

- Valid range: 2048-16384 bytes
- Optional

#### Authentication Configuration

For HTTP monitors that require authentication, specify the credentials reference under the `authentication` section:

```yaml
  authentication:
    secretRef: healthmonitor-secret
```

The referenced secret must contain `username` and `password` fields:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: healthmonitor-secret
  namespace: default
type: ako.vmware.com/basic-auth
data:
  username: <base64-encoded-username>
  password: <base64-encoded-password>
```

**Note**: The secret must exist in the same namespace as the HealthMonitor resource. When `authType` is specified in `httpMonitor`, this field is mandatory.

#### Federation Support

The `isFederated` field determines the replication scope of the health monitor object.

```yaml
  isFederated: false
```

- When set to `false` (default): The health monitor is visible only within the controller-cluster and its associated service engines
- When set to `true`: The health monitor is replicated across the federation

**Note**: This field is immutable once the HealthMonitor is created.

### Complete Examples

#### TCP Health Monitor with Half-Open Check

```yaml
apiVersion: ako.vmware.com/v1alpha1
kind: HealthMonitor
metadata:
  name: tcp-half-open-monitor
  namespace: default
spec:
  type: HEALTH_MONITOR_TCP
  sendInterval: 15
  receiveTimeout: 5
  successfulChecks: 2
  failedChecks: 3
  monitorPort: 80
  tcpMonitor:
    tcpHalfOpen: true
  isFederated: false
```

#### TCP Health Monitor with Request/Response

```yaml
apiVersion: ako.vmware.com/v1alpha1
kind: HealthMonitor
metadata:
  name: tcp-request-response-monitor
  namespace: default
spec:
  type: HEALTH_MONITOR_TCP
  sendInterval: 10
  receiveTimeout: 4
  successfulChecks: 2
  failedChecks: 2
  monitorPort: 8080
  tcpMonitor:
    tcpRequest: "GET /health HTTP/1.1"
    tcpResponse: "OK"
    maintenanceResponse: "Maintenance"
  isFederated: false
```

#### HTTP Health Monitor with Basic Authentication

```yaml
apiVersion: ako.vmware.com/v1alpha1
kind: HealthMonitor
metadata:
  name: http-auth-monitor
  namespace: default
spec:
  type: HEALTH_MONITOR_HTTP
  sendInterval: 10
  receiveTimeout: 5
  successfulChecks: 2
  failedChecks: 3
  monitorPort: 443
  authentication:
    secretRef: healthmonitor-secret
  httpMonitor:
    httpRequest: "GET /api/health HTTP/1.1"
    httpResponseCode:
      - HTTP_2XX
    httpResponse: "healthy"
    exactHttpRequest: false
    authType: AUTH_BASIC
    maintenanceCode:
      - 503
    maintenanceResponse: "Under Maintenance"
  isFederated: false
```

#### PING Health Monitor

```yaml
apiVersion: ako.vmware.com/v1alpha1
kind: HealthMonitor
metadata:
  name: ping-monitor
  namespace: default
spec:
  type: HEALTH_MONITOR_PING
  sendInterval: 10
  receiveTimeout: 4
  successfulChecks: 2
  failedChecks: 2
  isFederated: false
```

### Status Messages

The HealthMonitor CRD provides status information about the health monitor's state on the Avi Controller.

#### Programmed HealthMonitor

```bash
$ kubectl get healthmonitor
NAME                    AGE
my-health-monitor       5m
```

To view the detailed status:

```bash
$ kubectl get healthmonitor my-health-monitor -o yaml
```

The status section includes:

```yaml
status:
  uuid: "healthmonitor-uuid-1234"
  observedGeneration: 1
  lastUpdated: "2025-10-08T10:30:00Z"
  backendObjectName: "my-health-monitor"
  tenant: "admin"
  controller: "ako-crd-operator"
  dependencySum: 12345678
  conditions:
    - type: Programmed
      status: "True"
      reason: Created
      message: "HealthMonitor successfully created on Avi Controller"
      lastTransitionTime: "2025-10-08T10:30:00Z"
```

#### Status Fields

- **uuid** - Unique identifier of the health monitor object in the Avi Controller
- **observedGeneration** - The generation of the HealthMonitor resource that was most recently processed by the AKO CRD Operator
- **lastUpdated** - Timestamp when the health monitor object was last updated in the Avi Controller
- **backendObjectName** - Name of the health monitor object created in the Avi Controller, formatted as `ako-crd-operator-<cluster-name>--<sha1-hash>` where the hash is computed from `<namespace>-<name>`
- **tenant** - Avi tenant where the health monitor is created (e.g., "admin")
- **controller** - Name of the controller managing this resource (always "ako-crd-operator")
- **dependencySum** - Checksum of all dependencies for the health monitor, used to detect configuration changes

#### Status Conditions

The HealthMonitor status includes a `Programmed` condition that indicates whether the health monitor has been successfully processed:

**Condition Type: Programmed**

- **Status: True** - The health monitor has been successfully programmed on the Avi Controller
  - Reasons: `Created`, `Updated`
  
- **Status: False** - The health monitor failed to be programmed
  - Reasons: `CreationFailed`, `UpdateFailed`, `UUIDExtractionFailed`, `DeletionFailed`, `DeletionSkipped`

Example of a failed status:

```yaml
status:
  conditions:
    - type: Programmed
      status: "False"
      reason: CreationFailed
      message: "Failed to create HealthMonitor on Avi Controller: invalid monitor_port"
      lastTransitionTime: "2025-10-08T10:30:00Z"
```

### Using HealthMonitor with Gateway API

HealthMonitor CRDs can be used with Gateway API HTTPRoute resources by directly referencing them in HTTPRoute backendRefs filters using ExtensionRef:

```yaml
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: my-route
  namespace: default
spec:
  parentRefs:
    - name: my-gateway
  hostnames:
    - "app.example.com"
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: "/"
      backendRefs:
        - name: my-service
          port: 8080
          filters:
            - type: ExtensionRef
              extensionRef:
                group: ako.vmware.com
                kind: HealthMonitor
                name: my-health-monitor
```

### Using HealthMonitor with L4Rule

HealthMonitor CRDs can be referenced in L4Rule for Layer 4 LoadBalancer services using the `healthMonitorCrdRefs` field:

```yaml
apiVersion: ako.vmware.com/v1alpha2
kind: L4Rule
metadata:
  name: my-l4-rule
  namespace: default
spec:
  backendProperties:
    - port: 8080
      protocol: TCP
      healthMonitorCrdRefs:
        - my-health-monitor
        - my-backup-health-monitor
```
For more information about L4Rule, see the [L4Rule documentation](./l4rule.md).


### Conditions and Caveats

#### Immutable Fields

The following fields are immutable after the HealthMonitor is created:

- `type`: The health monitor type cannot be changed
- `monitorPort`: The monitor port cannot be modified
- `isFederated`: The federation setting cannot be changed

To modify these fields, the HealthMonitor must be deleted and recreated.

#### HealthMonitor Deletion

When a HealthMonitor is deleted:

- The health monitor object is removed from the Avi Controller if it is not being referenced by other objects
- The AKO CRD Operator retries the cleanup automatically

#### Validation Rules

The HealthMonitor CRD enforces the following validation rules:

1. **Authentication Requirement**: If `httpMonitor.authType` is set, `authentication.secretRef` must be configured.

2. **TCP Half-Open Constraints**: When `tcpMonitor.tcpHalfOpen` is true, the fields `tcpRequest`, `tcpResponse`, and `maintenanceResponse` must not be set.

3. **Exact HTTP Request with Auth**: If `httpMonitor.authType` is set, `httpMonitor.exactHttpRequest` must be false.

4. **Response Code Requirement**: For HTTP monitors, at least one value must be specified in `httpMonitor.httpResponseCode`.

5. **Maintenance Code Range**: Values in `httpMonitor.maintenanceCode` must be between 101 and 599.

#### Prerequisites

Before creating a HealthMonitor CRD:

1. The AKO CRD Operator must be installed and running in the cluster
2. For authenticated health checks, the secret containing credentials must be created first
3. Ensure the Avi Controller is accessible and properly configured

#### Namespace Considerations

- HealthMonitor resources can be created in any namespace
- When referenced by other resources (HTTPRoute/L4Rule), both the HealthMonitor and the referencing resource should typically be in the same namespace
- Secrets referenced in the `authentication` section must be in the same namespace as the HealthMonitor
- For HTTPRoute ExtensionRef: HealthMonitor must be in the same namespace as the HTTPRoute (ExtensionRef does not support cross-namespace references)
- For L4Rule: Uses simple string names (not object references with `kind` field) in the `healthMonitorCrdRefs` list. HealthMonitor CRDs must be in the same namespace as the L4Rule
- For RouteBackendExtension: HealthMonitor can also be referenced but only as AVIREF (reference to health monitor present on Avi Controller), not as CRD reference

#### Troubleshooting

If a HealthMonitor fails to be programmed:

1. Check the status conditions for detailed error messages:
   ```bash
   kubectl describe healthmonitor <name>
   ```

2. Verify the AKO CRD Operator logs:
   ```bash
   kubectl logs -n avi-system <ako-crd-operator-pod>
   ```

3. Ensure all referenced secrets exist and are properly formatted

4. Verify the Avi Controller connectivity and permissions

5. If a HealthMonitor is stuck in a terminating state:
   
   This typically occurs when the finalizer `healthmonitor.ako.vmware.com/finalizer` cannot be removed, usually due to:
   - AKO CRD Operator not running or unable to process the deletion   
   - The health monitor object on the Avi Controller is in use or cannot be deleted because it is referred by other objects
   
   To resolve:
   
   a. Check if the AKO CRD Operator is running:
      ```bash
      kubectl get pods -n avi-system | grep ako-crd-operator
      ```  
   
   b. Verify the health monitor on the Avi Controller and check if it's being referenced by other objects (Virtual Services, Pools, etc.). If referenced, remove those references first.
   
   c. If the operator is stuck and the Avi Controller object has been manually cleaned up, you can force remove the finalizer:
      ```bash
      kubectl patch healthmonitor <name> -n <namespace> -p '{"metadata":{"finalizers":[]}}' --type=merge
      ```
      
      **Warning**: Only use this as a last resort when confirmed the backend object is properly cleaned up on the Avi Controller. Removing the finalizer without proper cleanup may leave orphaned objects on the Avi Controller.


