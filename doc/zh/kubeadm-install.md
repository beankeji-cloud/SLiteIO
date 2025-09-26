# 使用kubeadm安装K8S

## 安装kubelet

```
cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kubelet kubeadm kubectl
EOF

# 设置SELinux为宽容模式
sudo setenforce 0
sudo sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

sudo yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes

sudo systemctl enable --now kubelet
```

## 启用CRI

```
# 注释掉disabled_plugins = ["cri"]
sudo vim /etc/containerd/config.toml

# 重启containerd
sudo systemctl restart containerd
```

## 安装K8S
```
# 将 "172.26.10.67" 替换为您机器的 IP 地址
sudo kubeadm init --ignore-preflight-errors Swap --apiserver-advertise-address=172.26.10.67 --pod-network-cidr=10.244.0.0/16
# 按照说明将 kubeconfig 文件复制到 $HOME/.kube/config 目录


kubectl create -f https://raw.githubusercontent.com/coreos/flannel/v0.22.0/Documentation/kube-flannel.yml

or
https://github.com/coreos/flannel/raw/master/Documentation/kube-flannel.yml
```

## 设置巨页
```
# 设置巨页
sudo bash -c "echo 256 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages"

# 重启kubelet
sudo systemctl restart kubelet

# 检查kubelet是否已识别巨页配置
kubectl get nodes -oyaml | grep hugepages-2Mi
```

## 检查状态

```
# 检查节点状态是否为就绪
$kubectl get node
NAME                                              STATUS   ROLES           AGE    VERSION
ip-172-26-10-67.ap-northeast-1.compute.internal   Ready    control-plane   139m   v1.27.4

# 所有 Pod 应处于运行状态
$kubectl get pods --all-namespaces
NAMESPACE      NAME                                                                      READY   STATUS    RESTARTS      AGE
kube-flannel   kube-flannel-ds-h5dcn                                                     1/1     Running   0             138m
kube-system    coredns-5d78c9869d-g4x5w                                                  1/1     Running   0             139m
kube-system    coredns-5d78c9869d-x8hmd                                                  1/1     Running   0             139m
kube-system    etcd-ip-172-26-10-67.ap-northeast-1.compute.internal                      1/1     Running   0             139m
kube-system    kube-apiserver-ip-172-26-10-67.ap-northeast-1.compute.internal            1/1     Running   0             139m
kube-system    kube-controller-manager-ip-172-26-10-67.ap-northeast-1.compute.internal   1/1     Running   0             139m
kube-system    kube-proxy-zs75b                                                          1/1     Running   0             139m
kube-system    kube-scheduler-ip-172-26-10-67.ap-northeast-1.compute.internal            1/1     Running   0             139m
```