# 构建 


## 构建 SLiteIO

### 环境要求

- linux or MacOS
- golang >= 1.17


### AMD64

```
# 构建 disk-controller 和 disk-agent
make controller

# 构建 CSI Driver
make csi

# 构建 scheduler plugin
make scheduler
```

### ARM64

以下命令将在 linux/arm64 平台构建相关组件（包括disk-agent和 CSI Driver），并将生成的镜像推送至 myregistry.com/SLiteIO/node-disk-controller:latest

```
PLATFORMS=linux/arm64 IMAGE_ORG=myregistry.com/SLiteIO TAG=latest make docker.buildx.agent
```


## 构建SPDK

建议在 CentOS 7.9 上构建 SPDK，该版本受到 SPDK 社区的良好支持

```
git clone https://github.com/spdk/spdk.git
cd spdk
git checkout v22.05
git submodule update --init

# 安装依赖项
./scripts/pkgdep.sh
yum install python3-pyelftools meson -y
 
# 如果不希望有此运行时依赖，请移除 libssl.so.1.1
yum remove openssl11-devel
mv /lib64/libssl.so.1.1 /lib64/libssl.so.1.1-backup

# 若需DPDK支持 NUMA，请安装 numa 库
yum install numactl-devel -y

# 编译
./configure
make
```
