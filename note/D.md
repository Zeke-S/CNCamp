# Docker 核心技术

## Dockerfile 常用指令

- **FROM:选择基础镜像**
- **LABELS:**
- **RUN:执行命令**
- **CMD:**
- **EXPOSE:发布端口**
- **ENV:设置环境变量**
- **ADD:源地址复制到目标路径**
- **ENTRYPOINT:定义可以执行的容器镜像入口命令**
- **VOLUME:**

## 传统服务vs微服务架构

微服务架构的缺点：分布式事务更加复杂（失败回滚）

微服务架构的优点：服务独立部署，易于持续集成和持续部署

## **微服务架构的设计和改造**

代码逻辑 + 业务逻辑 + 资源需求 三个角度出发去设计系统

分解原则：Size + Scope + Capabilities

## **常用的可以微服务化的组件**

用户和账户管理 - 用户管理 - 用户+用户组 - 业务逻辑出发

授权  - 基础权限 - 业务逻辑出发

系统配置 - 基础配置 - 业务逻辑出发

会话管理 - 代码逻辑和模块化功能出发

通知和通知服务 - 代码逻辑和模块化功能出发

照片、多媒体、元数据

## **微服务之间的通信方式**

点对点通信 + API网关通信

最佳实践方案，微服务之间可以采用点对点通信，对外开放平台采用API网关的形式完成统一的鉴权、调度

## docker 的常用命令

```shell
// 启动
docker run 
//	-it 交互
//	-d 后台运行
//	-p 端口映射
//	-v 磁盘挂载

// 启动已经终止的容器
docker start

// 停止容器
docker stop

// 查看容器进程
docker ps

// 查看容器的配置细节
docker inspect
// 运行状态
// 容器参数
// 镜像详情

// 打包、推送、拉去
docker build
docker push 
docker pull
```

## Namespace

资源的独立分配、进程隔离

```shell
// 进程数据结构
struct task_struct {
    // ...
    struct nsproxy *nsproxy;
    // ...
}

// Namespace 数据结构
struct nsproxy {
    atomic_t count;
    struct uts_namespace *uts_ns;
    struct ipc_namespace *ipc_ns;
    struct mnt_namespace *mnt_ns;
    struct pid_namespace *pid_ns_for_children;
    struct net *net_ns;
}
```

###  - Linux 对namespace 的操作方法：

1. clone

   在创建新进程的系统调用时，可以通过flags参数指定需要新建的Namespace

   // CLONE_NEWCGROUP / CLONE_NEWPIC / CLONE_NEWNET / CLONE_NEWNS / CLONE_NEWPID / CLONE_NEWUSER / CLONE_NEWUTS

2. setns

   该系统调用可以让调用进程加入到某个已经存在的Namespace

   int setns(int fd, int nstype)

3. unshare

   该系统调用可以让调用进程移动到新的namespace中去

   int unshare(int flags)

### -  Namespace的类型

| Namespace类型 |             隔离资源             | Kernel版本 |
| :-----------: | :------------------------------: | :--------: |
|      IPC      |  System V IPC 和 POSIX 消息队列  |   2.6.19   |
|    Network    | 网络设备、网络协议栈、网络端口等 |   2.6.29   |
|      PID      |               进程               |   2.6.14   |
|     Mount     |              挂载点              |   2.4.19   |
|      UTS      |          主机名和用户名          |   2.6.19   |
|      USR      |           用户和用户组           |    3.8     |

### - Namespace 的常用操作

1. 查看当前系统的namespace

   lsns -t <type>

2. 查看某个进程的namespace

   ls -la /proc/<pid>/ns/

   查看进程相关的namespace id

   proc是特殊的文件系统，查看当前主机运行的process的相关配置

3. 进入某个process的namespace运行命令

   nsenter -t <pid> -n ip addr

   // 进入某个进程<pid>的网络namespace运行ip addr

4. unshare调整一个进程的namespace

### - Namespace实验

```shell
// 启动一个进程 sleep 2分钟，但是让该进程进入新的网络namespace
sudo unshare -fn sleep 120
// 查看进程信息
sudo ps -ef | grep sleep
// 查看当前主机上所有的网络ns，并查看其对应的进程pid
sudo lsns -t net
// 查看对应进程pid的相关的ns, 观察其中的ns对应的编号，哪些是新的哪些是复用父进程的
sudo ls -al /proc/<pid>/ns
// 进入相关进程的net ns，查看具体的网络配置情况 
sudo nsenter -t <pid> -n ip addr
```

系统相关命令

```shell
// 查看网络配置
ip addr
ip -a
```

## CGroup - Control Group

对linux下用于对一个或者一组进程进行资源控制和监控的机制

资源子系统：CPU / MEM / 磁盘 / IO

CGroup在不同的资源管理子系统中以层级树（Hierarchy）的方式来组织管理，每个cgroup都可以包含其他的子cgroup，因此子cgroup能使用的资源除了受本cgroup配置的资源参数限制，还受父cgroup设置的资源限制。

```shell
struct task_struct {
    #ifdef CONFIG_CGPOUPS
    struct css_set_rcu *cgroups;
    struct list_head cg_list;
    #endif
}

struct css_set{
    /* 
     * cgroup subsystem states
     * set of subsystem states, one for eache subsystem.This array is immutable after creation apart from
     * init_css_set during subsystem registeration(at boot time).
     */
    struct cgroup_subsys_state *subsys[CGROUP_SUBSYS_COUNT];
}
```

cgroupp的层级结构: /sys/fs/cgroup

### - CGroup可控制的资源类型

- blkio: block io 子系统设置限制每个块设备的输入输出。例如：磁盘，光盘以及USB等等
- cpu: 该子系统使用调度程序为cgroup任务提供CPU的访问
- cpuacct: 产生cgroup任务的CPU资源报告
- cpuset: 如果是多核的CPU，该子系统会为cgroup任务分配单独的CPU和内存
- devices: 允许或者拒绝cgroup任务对设备的访问
- freezer: 暂停和恢复cgroup任务
- memory: 设置每个cgroup的内存限制以及产生内存资源报告
- net_cls: 标记每个网络包以供cgroup方便使用
- ns: 名称空间子系统
- pid: 进程标识子系统

### - CPU子系统

cpu.shares: 可出让的能获得CPU使用时间的相对值

cpu.cfs_period_us: 用来配置时间周期长度，单位us

cpu.cfs_quota_us: 用来配置当前cgroup在cfs_period_us时间内最多能使用CPU时间数，单位us

cpu.stat:  cgroup 内的进程使用的cpu时间统计

nr_periods: 经过cpu.cfs_period_us的时间周期数量

nr_throttled: 在经过的周期内，有多少次因为进程在指定的时间周期内用光了配额时间而收到限制。

throttled_time:  cgroup中的进程被限制使用CPU的总用时，单位是ns（纳秒）

### - Linux调度器

内核默认提供了5个调度器，Linux内核使用struct_sched_class来对调度器进行抽象

- Stop调度器，stop_sched_class: 优先级最高的调度器，可以抢占其他所有进程，不能被其他进程抢占；
- Deadline调度器， dl_sched_class: 使用红黑树，把进程按照绝对截止期限进行排序，选择最小进程进行调度运行；
- RT调度器，rt_sched_class: 实时调度器，为每个优先级维护一个队列；
- CFS调度器，cfs_sched_class: 完全公平调度器，采用完全公平的调度算法，引入虚拟运行时间概念；
- IDLE-Task调度器，idle_sched_class: 空闲调度器，每个CPU都会有一个idle线程，当没有其他进程可以调度时，调度运行idle线程；

### - CFS 调度器

相关命令

```shell
// 查看进程资源占用
top
```

### - CGroup实验 - CPU子系统

```shell
# 运行busyloop程序，启动两个线程占用cpu资源
# 查看进程的资源占用，并记录busyloop进程号
top
# 建立新的cgroup子系统
cd /sys/fs/cgroup/cpu
mkdir cpudemo
# 将新的进程<pid>写入cpudemo cgroup进程列表
sudo bash -c 'echo 24996 > cgroup.procs'
# 限制CPU的绝对时间, 只占用一个cpu
echo 10000 > cpu.cfs_quota_us
# 删除cpudemo
rmdir cpudemo
# 注意，直接采用rm -r的方式无法删除cgroup
```

### - Memory 子系统

memory.usage_in_bytes: 

memory.max_usage_in_bytes:

memory.limit_in_bytes:

memory.oom_control: 



