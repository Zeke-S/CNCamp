- **k8s**

  谷歌开源的容器集群管理系统，主要功能：

  1. 基于容器的应用部署、维护和滚动升级
  2. 负载均衡和服务发现
  3. 夸机器和跨地区的集群调度
  4. 自动伸缩
  5. 无状态服务和有状态服务
  6. 插件机制保证扩展性

- **命令式（Imperative）vs 声明式（Declarative）**

  命令式系统关注**如何做**

  声明式系统关注**做什么**

- **声明式系统规范**

  幂等型/面向对象

- **k8s核心对象**

  Node：计算节点的抽象，计算节点的资源抽象、健康状态

  Namespace：隔离资源的基本单位

  Pod：用来描述应用实例，基本的调度单元

  Service：如何将应用发布成服务，负载均衡和域名服务

- **Kubernetes主节点（Master Node）**

  API Server：

  Cluster Data Store：etcd

  Controller Manager：

  Scheduler：

- **Worker Node**

  Kubelet

  Kube-proxy

- **Raft协议**
























```shell
# kubeadm集群搭建


```

