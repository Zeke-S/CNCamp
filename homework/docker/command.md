**编写Dockerfile**

build/Dockerfile

注意事项：

	1. 需要设置交叉编译环境
 	2. 采用多段构建的方式优化镜像生成过程

**构建本地镜像**

```shell
cd build
docker build -t sunxiaocheninsun/httpserver:1.0 .
```

**推送镜像仓库**

```shell
docker push sunxiaocheninsun/httpserver:1.0
```

**启动镜像**

```shell
# docker pull sunxiaocheninsun/httpserver:1.0
docker run -i -p 8090:8090 --rm sunxiaocheninsun/httpserver:1.0 
```

**查看网络配置**

```shell
# 查看进程pid
sudo lsns -t net | grep httpserver
# 查看容器的IP配置 pid 55396 
sudo nsenter -t 55396 -n ip addr
```

