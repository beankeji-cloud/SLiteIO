# 部署与使用

## 环境说明

1. 操作系统
RHEL9
2. Kubernetes版本
v1.28.0+

## 部署方法

SLiteIO 需要部署在 K8S 环境中，并且有集群的管理员权限，需要至少一台可用节点。

### 准备工作

1. 安装内核模块（如果希望使用远程后端存储, 节点需要提前安装 nvme-tcp 内核模块）
```
modprobe nvme
modprobe nvme-tcp
```

2. 配置巨页
```
echo never > /sys/kernel/mm/transparent_hugepage/enabled
echo never > /sys/kernel/mm/transparent_hugepage/defrag
echo 256 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages
systemctl restart kubelet
modprobe dm-thin-pool
```

3. 配置节点是否支持精简卷
```
## 替换[node-name]，支持精简卷
kubectl label node [node-name] lite.io/thin=true
## 替换[node-name]，不支持精简卷
kubectl label node [node-name] lite.io/thin=false
```

4. 修改块设备名称
1）hack/deploy/lvm/050-configmap.yaml文件
```
## 范例（共一处，本例块设备名称是/dev/sdb）
devicePath: /dev/sdb
```

### 部署SLiteIO

在项目根目录下，使用以下命令部署组件

1. 安装 SLiteIO 组件
```
kubectl create -f hack/deploy/base
```

2. 安装 LVM 数据引擎配置映射
```
kubectl create -f hack/deploy/lvm
```

### 验证部署SLiteIO

1. 检查容器是否运行正常
```
## 范例
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

2. 检查 StorageClass 是否正确创建
```
## 范例
# kubectl get sc | grep nvmf
## 采用默认配置，无强制或优先的节点位置策略，自动选择存储节点。
antstor-nvmf                     antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 卷必须创建在与 Pod 调度到的节点相同的本地存储上。
antstor-nvmf-mustlocal           antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 卷必须创建在与 Pod 调度节点不同的远程存储节点上。
antstor-nvmf-mustremote          antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 优先在 Pod 所在节点的本地存储创建卷，若本地存储不可用，再选择远程存储节点。
antstor-nvmf-preferlocal         antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 优先在远程存储节点创建卷，若远程存储不可用，再选择本地存储节点。
antstor-nvmf-preferremote        antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 采用精简置备模式，无强制或优先的节点位置策略，自动选择存储节点。
antstor-nvmf-thin                antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 采用精简置备模式，卷必须创建在与 Pod 调度到的节点相同的本地存储上。
antstor-nvmf-thin-mustlocal      antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 采用精简置备模式，卷必须创建在与 Pod 调度节点不同的远程存储节点上。
antstor-nvmf-thin-mustremote     antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 采用精简置备模式，优先在 Pod 所在节点的本地存储创建卷，若本地存储不可用，再选择远程存储节点。
antstor-nvmf-thin-preferlocal    antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 采用精简置备模式，优先在远程存储节点创建卷，若远程存储不可用，再选择本地存储节点。
antstor-nvmf-thin-preferremote   antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 基于 LVM 卷组的 nvmf 存储类
antstor-nvmf-vg                  antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
## 采用精简置备模式，基于 LVM 卷组的 nvmf 存储类
antstor-nvmf-vg-thin             antstor.csi.alipay.com   Delete          WaitForFirstConsumer   true                   52m
```

## 使用方法

1. 使用 hack/deploy/example/pod.yaml 文件创建pod
```
kubectl create -f hack/deploy/example/pod.yaml
```

2. 检查创建结果
```
## 范例
# kubectl -n obnvmf get pod |grep test-pod
test-pod                                  1/1     Running   0             57m
# kubectl -n obnvmf get antstorvolume
NAME                                       UUID                                   SIZE           THIN   TARGETID   HOST_IP         STATUS   AGE
pvc-b4cee247-34d6-4a56-882a-213de10fae1b   ec82f398-f3fd-4736-ac12-04e3d772d0bf   161061273600   true   k8s-179    10.243.33.179   ready    172m
```

3. 检查卷健康
```
# 范例（更换[node-name]名称）
# kubectl get --raw /api/v1/nodes/[node-name]/proxy/metrics | grep "health_status_abnormal"
## 卷健康状态
kubelet_volume_stats_health_status_abnormal{namespace="obnvmf",persistentvolumeclaim="pvc-obnvmf-test"} 0
```

4. 查看Prometheus数据
```
## 范例（更换[node-ip]地址）
# curl http://[node-ip]:6010/metrics
## 本地卷数量
liteio_node_local_volume_num{node="k8s-178",thin="true"} 0
## 远程卷数量
liteio_node_remote_volume_num{node="k8s-178",thin="true"} 0
## 精简池已分配的容量
liteio_thinpool_allocated_bytes{node="k8s-178",thin="true"} 0
## 精简池data空间占用率
liteio_thinpool_data_percent{node="k8s-178",thin="true"} 0
## 精简池metadata空间占用率
liteio_thinpool_metadata_percent{node="k8s-178",thin="true"} 0.1041
## 精简池超分比
liteio_thinpool_overprovision_ratio{node="k8s-178",thin="true"} 2
## 剩余容量
liteio_vg_available_bytes{node="k8s-178",thin="true"} 5.3659828224e+11
## 总容量
liteio_vg_size_bytes{node="k8s-178",thin="true"} 5.3659828224e+11
```