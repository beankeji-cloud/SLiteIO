# Practical Implementation of a Customized Version Based on Ant Group's Open-Source LiteIO

Based on the open-source repository of LiteIO (https://github.com/eosphoros-ai/liteio), we have developed a storage solution tailored for containerization transformation in data centers, particularly for environments built on traditional IT infrastructure such as standard networking and SATA-SSDs. This solution proves especially critical when containerizing databases. Our customized version has been applied to MySQL containerization, validated through systematic testing, and is being gradually rolled out in enterprise production environments.

# What Is SLiteIO?

LiteIO is originally designed to provide a high-performance and scalable cloud-native block storage service. However, in reality, many data centers have a large number of existing/repurposed devices with 10 Gigabit Ethernets and standard SSDs. These enterprises often lack the budget to purchase NVMe SSDs or IB network devices. Therefore, SLiteIO has been repositioned to provide a container storage solution for these enterprises.

## **Design Background**

To reduce costs and improve efficiency, many enterprises are replacing traditional virtualization with containers. Containers offer significant advantages in stateless applications, but when it comes to stateful applications,such as databases,determining the appropriate storage solution becomes a challenge. 

Choosing local storage gives rise to the following two issues:

- **Uneven Utilization**: I/O-intensive and compute-intensive workload vary, leading to scenarios where one machine may be fully utilized for computation while storage remains underutilized, or vice versa. Moreover, attaining a globally optimal solution through scheduling is a considerable challenge.
- **Poor Scalability**: When the storage is insufficient and the storage needs to be scale up, it becomes essential to migrate to a server with a larger storage capacity, which takes a long time to copy data.



Traditional distributed storage represents a decent solution, but within the domain of databases, it introduces several problems:

- **Ascension In Replication Count**(Cost): The advantage of distributed storage lies in the pooling storage through erasure coding (EC) and multi-replica techniques, which offer robust protection against single hardware failures. However, under this architecture, the application of EC and multi-replica results in a replication factor greater than 1 for each data segment, usually between 1.375 and 2. As an important component of business services, databases often necessitate geo-redundancy and disaster recovery across different availability zones (AZ) at the upper layer, with a backup replication already existing in another AZ. The total number of data replicas is set to rise.
- **Large Explosion Radius**(Stability): Distributed storage typically features a centralized metadata layer, when subject to failure, can lead to global exceptions.

## **Design Ideas**

The design of Sliteio strives to simplify the data path. This not only reduces latency, but also avoids stability risks caused by complex path management. Sliteio directly adopts LVM as its data engine, incorporates thin provisioning management, and can export remote volumes via SPDK. Compute nodes can access storage nodes remotely through NVMe over TCP. The peer-to-peer design, combined with Kubernetes' scheduling control, effectively mitigates the impact of a single hardware failure on services.

![image](../image/architecture-en.jpg)


## **Cost Management**

Based on SLiteIO, unallocated storage within servers can be dynamically distributed to remote compute nodes on demand. The globally coordinated scheduling pools global storage resources, thereby enhancing overall storage efficiency.

For example, there are two types of servers: compute-intensive (96C4T) and storage-intensive(64C20T). Assuming that the CPU of the storage model has been allocated, with 5TB disk space left, while the compute model still has available CPU but no disks to allocate. Using LiteIO, it is possible to combine the CPU of the compute model with the remaining disk space of the storage model into a new container to provide services, thus enhancing the efficiency of both computational and storage capacities.

## General Storage And Computing Separation

SLiteIO is a general storage service technology, which acts on storage logic volumes, in conjunction with K8S, the storage perceived by upper-layer containers or applications is indistinguishable from that of local disks. Whether it is direct read/write to block devices (bdev) or formatting the block devices into any file system, no modifications are required from upper-layer services. Databases such as OceanBase, MySQL, and PostgreSQL written in Java, Python, or Go can utilize it as a regular disk.

## **Serverless**

SLiteIO's general storage and computing separation capability simplifies scaling dramatically. With the perception and scheduling system, deploying a MySQL instance inherently gains serverless capability. When the computing power of MySQL is insufficient, one can rapidly achieve scale-up by attaching MySQL storage to a more powerful container via SLiteIO. When storage space of MySQL is insufficient, simply mounting an additional disk from another storage node allows for expansion without data loss.

# **Technical Features**

## **Transport Protocol**

SLiteIO uses NVMe-oF for cross-node data transmission. Even without NVMe SSDs and RDMA networks, it can deliver impressive performance using standard SSDs and 10Gb Ethernet cards.

## **Simplified IO Pipeline**

In the traditional distributed storage architectures, a write I/O operation involves three steps: querying metadata, writing metadata, and writing multiple data replicas, which require numerous network interactions. In the SLiteIO architecture, due to one-to-one mapping between frontend PVC and backend volumes, no additional rootserver or metaserver is required to manage global metadata and data-block. The I/O path requires only a single network interaction, eliminating the latency and amplification issues associated with multi-replica writes,which makes SLiteIO lower and more stable I/O latency.

## **Multi-Disk**

With the one-to-one volume model, when a node's storage capacity is nearly full, some resource fragments  inevitably remain. SLiteIO can aggregate these fragments into a single volume for applications utilization. This approach introduces a higher failure rate issue: if any node providing the fragments fails, the volume becomes unavailable.Therefore, it is recommended that critical services, especially core databases, avoid using this method.

## **Thin Provisioning**

SLiteIO also offers Thin Provisioning capability, which enable over-provisioning of storage. In practice, combining a reasonable over-provisioning ratio with appropriate scheduling policies can significantly improve storage utilization.

## **Capacity Expansion**

SLiteIO supports capacity expansion.When used in conjunction with thin provisioning, the specific threshold of the over-allocation ratio must be considered to prevent risks associated with excessive capacity over-allocation.

# **Practice**

During our containerization of MySQL, we implemented a storage solution based on SLiteIO. After containerization, database read and write performance significantly outperformed local storage in virtualized environments, with minimal latency difference between remote and local access.

