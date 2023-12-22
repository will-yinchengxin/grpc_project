## RocketMQ
compose 文件
```dockerfile
version: '3.5'
services:
  rmqnamesrv:
    image: rocketmqinc/rocketmq
    container_name: rmqnamesrv
    restart: always
    ports:
      - 9876:9876
    environment:
    #内存分配
      JAVA_OPT_EXT: "-server -Xms1g -Xmx1g"
    volumes:
      - ./logs:/root/logs
    command: sh mqnamesrv
    networks:
      rmq:
        aliases:
          - rmqnamesrv
          
  rmqbroker:
    image: dyrnq/rocketmq:4.8.0
    container_name: rmqbroker
    restart: always
    depends_on:
      - rmqnamesrv
    ports:
      - 10909:10909
      - 10911:10911
    volumes:
      - ./logs:/root/logs
      - ./store:/root/store
      - ./conf/broker.conf:/opt/rocketmq/conf/broker.conf
    command: sh mqbroker  -c ../conf/broker.conf
    environment:
      NAMESRV_ADDR: "rmqnamesrv:9876"
      JAVA_OPT_EXT: "-server -Xms1g -Xmx1g -Xmn1g"
    networks:
      rmq:
        aliases:
          - rmqbroker
          
  rmqconsole:
    image: apacherocketmq/rocketmq-dashboard
    container_name: dashboard
    restart: always
    ports:
      - 8080:8080
    depends_on:
      - rmqnamesrv
    environment:
      JAVA_OPTS: "-Drocketmq.namesrv.addr=rmqnamesrv:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false"
    networks:
      rmq:
        aliases:
          - rmqconsole
          
networks:
  rmq:
    name: rmq
    driver: bridge
````
broker.conf文件
````
brokerClusterName = DefaultCluster
brokerName = broker-a
brokerId = 0
deleteWhen = 04
fileReservedTime = 48
brokerRole = ASYNC_MASTER
flushDiskType = ASYNC_FLUSH
````
启动服务
```shell
docker-compose up -d
````
访问
````
http://localhost:8080/#/consumer
````
