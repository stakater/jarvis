
##### Node Conditions

- Apart from the default conditions supported by native k8s, `NodeProblemDetector` provides additional NodeConditions based 
on the `ProblemDaemon` enabled with NodeProblemDetector.

- K8s Native NodeConditions:
  
  | ConditionType      | Condition Status   |Effect        | Key      |
  | ------------------ | ------------------ | ------------ | -------- |
  |Ready               |True                | -            | |
  |                    |False               | NoExecute    | node.kubernetes.io/not-ready           |
  |                    |Unknown             | NoExecute    | node.kubernetes.io/unreachable         |
  |OutOfDisk           |True                | NoSchedule   | node.kubernetes.io/out-of-disk         |
  |                    |False               | -            | |
  |                    |Unknown             | -            | |
  |MemoryPressure      |True                | NoSchedule   | node.kubernetes.io/memory-pressure     |
  |                    |False               | -            | |
  |                    |Unknown             | -            | |
  |DiskPressure        |True                | NoSchedule   | node.kubernetes.io/disk-pressure       |
  |                    |False               | -            | |
  |                    |Unknown             | -            | |
  |NetworkUnavailable  |True                | NoSchedule   | node.kubernetes.io/network-unavailable |
  |                    |False               | -            | |
  |                    |Unknown             | -            | |
  |PIDPressure         |True                | NoSchedule   | node.kubernetes.io/pid-pressure        |
  |                    |False               | -            | |
  |                    |Unknown             | -            | |


- `NodeProblemDetector`(NPD) supported NodeConditions:
  - NPD only patches Nodes with conditions, it doesn’t apply taints on Nodes.
We will have to decide on the effect and taint against the following conditions supported by NPD 

- `ntp-custom-plugin-monitor`
  -   | ConditionType      | Condition Status   |Effect        | Key      |
      | ------------------ | ------------------ | ------------ | -------- |
      |NTPProblem          |True                |              | |
      |                    |False               |              |          |
      |                    |Unknown             |              |          |

- `docker-monitor`
  -   | ConditionType        | Condition Status   |Effect        | Key      |
      | ------------------   | ------------------ | ------------ | -------- |
      |CorruptDockerOverlay2 |True                |              | |
      |                      |False               |              |          |
      |                      |Unknown             |              |          |

- `Health-checker-containerd, docker`
  -   | ConditionType            | Condition Status   |Effect        | Key      |
      | ------------------       | ------------------ | ------------ | -------- |
      |ContainerRuntimeUnhealthy |True                |              | |
      |                          |False               |              |          |
      |                          |Unknown             |              |          |

- `Health-checker-kubelet`
  -   | ConditionType            | Condition Status   |Effect        | Key      |
      | ------------------       | ------------------ | ------------ | -------- |
      |KubeletUnhealthy          |True                |              | |
      |                          |False               |              |          |
      |                          |Unknown             |              |          |

- `kernel-monitor`
  -   | ConditionType            | Condition Status   |Effect        | Key      |
      | ------------------       | ------------------ | ------------ | -------- |
      |KernelDeadlock            |True                |              | |
      |                          |False               |              |          |
      |                          |Unknown             |              |          |
      |ReadonlyFilesystem        |True                |              | |
      |                          |False               |              |          |
      |                          |Unknown             |              |          |
      |FrequentUnregisterNetDevice        |True                |              | |
      |                                   |False               |              |          |
      |                                   |Unknown             |              |          |

- `systemd-monitor`
  -   | ConditionType            | Condition Status   |Effect        | Key      |
      | ------------------       | ------------------ | ------------ | -------- |
      |FrequentKubeletRestart    |True                |              | |
      |                          |False               |              |          |
      |                          |Unknown             |              |          |
      |FrequentDockerRestart     |True                |              | |
      |                          |False               |              |          |
      |                          |Unknown             |              |          |
      |FrequentContainerdRestart          |True                |              | |
      |                                   |False               |              |          |
      |                                   |Unknown             |              |          |   
  

---

##### One node can have not good several conditions at one time. We need to new condition types that could be combination multiple condition types.
- In certain case where a single condition is not sufficient to mark a node as unhealthy we can support a new type as `ConditionSet`.

##### Approach 1

##### ConditionSet as ConfigMap

- A `ConditionSet` would be combination of 1 or more conditions.
- Sample ConfigMap
  ```yaml
    type: ConditionSets
    conditionSets:
    - type: KubeletContainerRuntimeUnhealthy
      effect: NoExecute
      taintKey: node.stakater.com/KubeletContainerRuntimeUnhealthy
      conditions:
        - ConditionType: KubeletUnhealthy
          conditionStatus: true
        - ConditionType: ContainerRuntimeUnhealthy
          conditionStatus: Unknown
    - type: KernelDeadlock
      effect: NoExecute
      taintKey: node.stakater.com/KernelDeadlock
      conditions:
        - ConditionType: KernelDeadlock
          conditionStatus: true
  ```
- If, a Node's conditions matches any of the ConditionSet then corresponding effect would be applied. 
- In any case Node's conditions matches multiple ConditionSets then higher level of effect(`NoExecute > NoSchedule > PreferNoSchedule`)  would be applied.


##### Approach 2

##### ConditionSet as a new CustomResource type

-
    ```yaml
    apiVersion: autohealer.resourcestack.com/v1alpha1
    kind: ConditionSet
    metadata:
      name: conditionset-1
    spec:
      type: KubeletContainerRuntimeUnhealthy
      effect: NoExecute
      taintKey: node.stakater.com/KubeletContainerRuntimeUnhealthy
      conditions:
        KubeletUnhealthy:
          status: true
        ContainerRuntimeUnhealthy:
          status: unknown
    ---
    apiVersion: autohealer.resourcestack.com/v1alpha1
    kind: ConditionSet
    metadata:
      name: conditionset-2
    spec:
      type: KernelDeadlock
      effect: NoExecute
      taintKey: node.stakater.com/KernelDeadlock
      conditions:
        KernelDeadlock:
          status: true
    ```
- The advantage of having a `ConditionSet` as a new resource type is we can apply [Validation Webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#validatingadmissionwebhook), 
  which can validate each `ConditionSet`.
  - For example, as we are giving the power to configure conditions and its effect.
    If we had set `effect` to `NoExecute` for condition type `OutOfDisk` with condition status as `false` instead of `true`, 
    then it will start evicting all healthy nodes in the cluster and could take down the whole cluster if not noticed.
- So by having a validation webhook, we can make sure that certain `condition` doesn't get configured which could impact the cluster.
- Another advantage could be avoiding duplicate condition entries, or handling certain not supported condition type.
