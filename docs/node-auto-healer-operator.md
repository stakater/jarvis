

### Node Auto Healer Operator (NAHO)

##### NodeAutoHealer

- `NodeAutoHealer` is an object or CR of NAHO, which provides how NAHO can be enabled to provide node auto healing.
- 
  ```yaml
    apiVersion: autohealer.resourcestack.com/v1alpha1
    kind: NodeAutoHealer
    metadata:
      name: node-auto-healer-1
    spec:
      disableAutoHealing: false
      nodeSelector:
        matchLabels:
          type: small 
        matchExpressions:
          - {key: tier, operator: In, values: [cache]}
          - {key: environment, operator: NotIn, values: [dev]}
      noScheduleThresholdLimit: 30mins     
      parallelHealing:
        enable: true
        maxAllowedNodesToDrain: 20%
      forceDelete: true
    status:
      state: active
      disabledAt: {LastTime when disabled}
  ```
- **kind**: NodeAutoHealer
- **spec**:  
  - **disableAutoHealing**: A boolean flag, to define change state of NAHO. Setting it to `true` will pause auto healing.
  - **nodeSelector**: If provided, nodes can be filtered using `matchLabels` or `matchExpression` and only those nodes would be considered for auto healing.
    - If not provided all nodes within the cluster would be considered for auto healing
  - **noScheduleThresholdLimit**: Threshold time limit on how long a node can remain in `NoSchedule` state. And if it exceeds the threshold we can apply `NoExecute` taint, which will evict all pods and later the Node will be drained.
  - **parallelHealing**: If enabled, then multiple nodes can be drained in parallel. 
    - **maxAllowedNodesToDrain**: Maximum number of nodes that can remain under drained state at any given time. 
      - Value can be an absolute number (ex: 5), or a percentage of total nodes at that particular moment (ex: 10%).
  - **forceDelete**: If set true, it will delete the node even if it fails to drain the node.
- **status**:
  - **state**: {Active|Paused}, represents whether auto healer is active or in pause state.
  - **disabledAt**: Last datetime, when auto healing was disabled

---

##### HealedNode

- `HealedNode` is an object or CR of NAHO, which will be created when a node will match certain `ConditionSet` and requires recovery.
- Using this CR we can monitor the progress of node's healing process. (drain node -> delete machine -> monitor, creation of replacement node)  
-
  ```yaml
  apiVersion: autohealer.resourcestack.com/v1alpha1
  kind: HealedNode
  metadata:
    name: kind-control-plane
  spec:
    nodeDetails:
      nodeName: kind-control-plane
      taints: 
        - effect: NoSchedule
          key: key1
          value: value1  
        - effect: NoExecute
          key: key1
          value: value1
      conditions:
      - lastHeartbeatTime: "2021-04-25T11:51:06Z"
        lastTransitionTime: "2021-04-25T11:36:05Z"
        message: Kubelet never posted node status.
        reason: NodeStatusNeverUpdated
        status: "Unknown"
        type: Ready
      addresses:
      - address: 172.18.0.2
        type: InternalIP
      - address: kind-control-plane
        type: Hostname
      nodeSystemInfo:
        architecture: amd64
        bootID: 3b622bbf-a04c-4a50-81d9-7afb89502684
        containerRuntimeVersion: containerd://1.4.0-106-gce4439a8
        kernelVersion: 5.8.0-50-generic
        kubeProxyVersion: v1.20.2
        kubeletVersion: v1.20.2
        machineID: bed729392962410587918db70d475183
        operatingSystem: linux
        osImage: Ubuntu 20.10
        systemUUID: 1e3f9f51-2c77-473d-a173-3e095e6e652c
    matchedConditionSets:
      - name: NetworkUnavailable
    appliedAction: {drainNode|deleteNode}
  status:
    currentState: {draining|drained|deleting|deleted|recovering|recovered}
    lastStateChangeTime:
    isHealingProcessStable: true   
  ```
  
- **kind**: HealedNode
- **spec**:
  - **nodeDetails** - Target node's related details such as its name, address, system info, taints, conditions etc.
  - **matchedConditionSets** - Matched ConditionSets that shows the cause of a node being unhealthy.
  - **appliedAction** - Once a node becomes unhealthy, it should be drained first and then the node should be deleted.
    - So we will support two types of actions `drainNode` and `deleteNode`.
- **status**:
  - **currentState**: 
    - Once action type `drainNode` would be applied, the node will go through these 2 states:
      - `draining`: state representing execution of pod evictions process.
      - `drained`: state representing completion of pod eviction.
    - After completion of draining process, the Machine object would be deleted and would through these 2 states:    
      - `deleting`: state representing deleting of machine object, once `HealedNode`'s `currentState` becomes `drained`.
      - `deleted`: state representing completion of deletion of machine object.
    - Once the machine object would be deleted, it would go through these 2 states:
      - `recovering`: Before deleting the machine object, we will store count of total number of nodes.
        - Next, once delete is complete, it will keep checking if total number of previous nodes equals total number of current nodes.
      - `recovered`: Once the count of total previous nodes becomes same as current total nodes, the state would become `recovered`.  
    