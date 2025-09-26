# Deployment and Usage
The Guide helps you to setup a local K8S cluster and deploy LiteIO in it. It is only for testing purpose.

## Preparing Node

1. OS：RHEL9
2. K8S：v1.28.0+
3. CPU/Memory: 4C8G at least

## How to Deploy

SLiteIO needs to be deployed in a Kubernetes environment and requires cluster administrator privileges, with at least one available node.

### Install dependencies

1. Install the kernel modules （Nodes should have the nvme-tcp kernel module installed in advance if remote backend storage will be used.）
```
modprobe nvme
modprobe nvme-tcp
```

2. Configure Huge Pages
```
echo never > /sys/kernel/mm/transparent_hugepage/enabled
echo never > /sys/kernel/mm/transparent_hugepage/defrag
echo 256 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages
systemctl restart kubelet
modprobe dm-thin-pool
```

3. Configure node support for thin provisioning capability
```
## Replace[node-name]，Support for Thin-Provisioned Volumes
kubectl label node [node-name] lite.io/thin=true
## Replace[node-name]，No support for Thin-Provisioned Volumes
kubectl label node [node-name] lite.io/thin=false
```

4. Rename block device
1）hack/deploy/lvm/050-configmap.yaml
```
## Example (There is one instance in total. In this example, the block device name is /dev/sdb)
devicePath: /dev/sdb
```

### Install SLiteIO

Navigate to the project root and run the following command to deploy the component.

1. Install SLiteIO
```
kubectl create -f hack/deploy/base
```

2. Install the LVM with mapping configuration
```
kubectl create -f hack/deploy/lvm
```

### Verify SLiteIO

1. Check the Status of POD
```
## Example
# kubectl -n obnvmf get pods
NAME                                      READY   STATUS    RESTARTS   AGE
csi-antstor-controller-5b5645b985-kbd5r   6/6     Running   0          18h
node-disk-controller-5878d88d5c-x67t9     1/1     Running   0          18h
nvmf-tgt-2wjfg                            1/1     Running   0          18h
nvmf-tgt-6c67s                            1/1     Running   0          18h
nvmf-tgt-qzjd5                            1/1     Running   0          18h
obnvmf-csi-node-cgknt                     3/3     Running   0          18h
obnvmf-csi-node-hzm5f                     3/3     Running   0          18h
obnvmf-csi-node-z7w9v                     3/3     Running   0          18h
obnvmf-disk-agent-app-2hvl9               1/1     Running   0          18h
obnvmf-disk-agent-thin-app-bzcp4          1/1     Running   0          18h
obnvmf-disk-agent-thin-app-xlw5q          1/1     Running   0          18h
```

2. Check the StorageClass
```
## Example
# kubectl get sc | grep nvmf
## Uses the default configuration with no mandatory or preferred node placement strategy; storage nodes are selected automatically.
antstor-nvmf                     antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Volumes must be created on local storage attached to the same node where the Pod is scheduled.
antstor-nvmf-mustlocal           antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Volumes must be created on remote storage nodes different from the node where the Pod is scheduled.
antstor-nvmf-mustremote          antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Preferentially create volumes on the local storage of the Pod's node; if local storage is unavailable, then select a remote storage node.
antstor-nvmf-preferlocal         antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Preferentially create volumes on a remote storage node; if remote storage is unavailable, then use local storage.
antstor-nvmf-preferremote        antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Uses thin provisioning mode with no mandatory or preferred node placement strategy; storage nodes are selected automatically.
antstor-nvmf-thin                antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Uses thin provisioning mode. Volumes must be created on local storage attached to the same node where the Pod is scheduled.
antstor-nvmf-thin-mustlocal      antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Uses thin provisioning mode. Volumes must be created on a remote storage node, different from the node where the Pod is scheduled.
antstor-nvmf-thin-mustremote     antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Uses thin provisioning mode. Preferentially create volumes on the local storage of the Pod's node; if local storage is unavailable, then select a remote storage node.
antstor-nvmf-thin-preferlocal    antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## Uses thin provisioning mode. Preferentially create volumes on a remote storage node; if remote storage is unavailable, then use local storage.
antstor-nvmf-thin-preferremote   antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## NVMe-oF StorageClass based on LVM Volume Groups.
antstor-nvmf-vg                  antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## NVMe-oF (NVMF) StorageClass that uses thin provisioning and is based on LVM Volume Groups.
antstor-nvmf-vg-thin             antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
```

## How to Use

1.Use hack/deploy/example/pod.yaml to create the Pod
```
kubectl create -f hack/deploy/example/pod.yaml
```

2. Check the Pod
```
## Example
# kubectl -n obnvmf get pod |grep test-pod
test-pod                                  1/1     Running   0             57m
# kubectl -n obnvmf get antstorvolume
NAME                                       UUID                                   SIZE           THIN   TARGETID   HOST_IP         STATUS   AGE
pvc-b4cee247-34d6-4a56-882a-213de10fae1b   ec82f398-f3fd-4736-ac12-04e3d772d0bf   161061273600   true   k8s-179    10.243.33.179   ready    172m
```

3. Check the Volumes
```
# Example（Rename[node-name]）
# kubectl get --raw /api/v1/nodes/[node-name]/proxy/metrics | grep "health_status_abnormal"
## Status of Volumes
kubelet_volume_stats_health_status_abnormal{namespace="obnvmf",persistentvolumeclaim="pvc-obnvmf-test"} 0
```

4. Check Prometheus Data
```
## Example（Replace[node-ip]）
# curl http://[node-ip]:6010/metrics
## Number of Local Volumes
liteio_node_local_volume_num{node="k8s-178",thin="true"} 0
## Number of Remote Volumes
liteio_node_remote_volume_num{node="k8s-178",thin="true"} 0
## Thin Pool Allocated Capacity
liteio_thinpool_allocated_bytes{node="k8s-178",thin="true"} 0
## Thin Pool Data Space Utilization
liteio_thinpool_data_percent{node="k8s-178",thin="true"} 0
## Thin Pool Metadata Space Utilization
liteio_thinpool_metadata_percent{node="k8s-178",thin="true"} 0.1041
## Thin Pool Overcommitment Ratio
liteio_thinpool_overprovision_ratio{node="k8s-178",thin="true"} 2
## Remaining Capacity
liteio_vg_available_bytes{node="k8s-178",thin="true"} 5.3659828224e+11
##Total Capacity
liteio_vg_size_bytes{node="k8s-178",thin="true"} 5.3659828224e+11
```