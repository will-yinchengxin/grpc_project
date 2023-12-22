### [gorm](https://gorm.io/zh_CN/docs/models.html)

### [nacos](https://nacos.io/zh-cn/docs/quick-start.html)
```shell
Clone 项目
git clone https://github.com/nacos-group/nacos-docker.git

cd nacos-docker

单机模式 Derby
docker-compose -f example/standalone-derby.yaml up
````
验证
```shell
docker ps -a
CONTAINER ID   IMAGE                       COMMAND                  CREATED        STATUS        PORTS                                                                                  NAMES
e8779eb3e47f   prom/prometheus:latest      "/bin/prometheus --c…"   20 hours ago   Up 4 hours    0.0.0.0:9090->9090/tcp, :::9090->9090/tcp                                              prometheus
b8c371701761   grafana/grafana:latest      "/run.sh"                20 hours ago   Up 4 hours    0.0.0.0:3000->3000/tcp, :::3000->3000/tcp                                              grafana
dbf1a75174ff   nacos/nacos-server:v2.2.3   "bin/docker-startup.…"   20 hours ago   Up 4 hours    0.0.0.0:8848->8848/tcp, :::8848->8848/tcp, 0.0.0.0:9848->9848/tcp, :::9848->9848/tcp   nacos-standalone
32725b379818   tancloud/hertzbeat          "./bin/entrypoint.sh"    2 months ago   Up 20 hours   22/tcp, 1158/tcp, 0.0.0.0:1157->1157/tcp, :::1157->1157/tcp                            hertzbeat
````
访问: http://172.16.27.95:8848/nacos

### [consul](https://www.consul.io/)
#### docker 安装
```shell
// https://hub.docker.com/_/consul
docker pull consul:1.15

docker run -d -p 8500:8500 -e CONSUL_BIND_INTERFACE='eth0' --name=consul1 consul:1.15 agent -server -bootstrap -ui -client='0.0.0.0' 
```

#### mac 上安装: [下载地址](https://releases.hashicorp.com/consul)
解压到一个目录后, 执行如下命令
```shell
#!/bin/bash
./consul agent -dev -bind=127.0.0.1 -client=0.0.0.0 -data-dir=./data
````