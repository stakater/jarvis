# JARVIS
Machine auto healer!



## Problem



For a kubernetes cluster to remain in a healthy state, all the nodes should remain in a healthy, running state.


## Solution



- Machine auto healer operator will always try to keep the nodes (machines)
  in your cluster in a healthy, running state.
- It will perform periodic checks on the health state of each node (machine) in your cluster.
- If a node (machine) fails consecutive health checks over an extended time period,
  it will initiate a repair process for that node (machine).



![](./docs/images/machine_auto_healer.png)

##### Node Conditions
