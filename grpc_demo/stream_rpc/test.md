```shell
cd /grpc_test/stream_rpc/server/main.go  && go run main.go

cd /grpc_test/stream_rpc/client/main.go && go run main.go
````

---


````
server: 
2023/11/15 17:17:30 Start Server 8090 !!!
2023/11/15 17:17:33 Got Client message [First msg] at point (0, 1)
2023/11/15 17:17:33 Got Client message [Second msg] at point (0, 2)
2023/11/15 17:17:33 Got Client message [Third msg] at point (0, 3)
2023/11/15 17:17:33 Got Client message [Fourth msg] at point (0, 4)
2023/11/15 17:17:33 Got Client message [Fifth msg] at point (0, 5)
2023/11/15 17:17:33 Got Client message [Sixth msg] at point (0, 6)
2023/11/15 17:17:33 All Chat Use 321000 Second

client: 
2023/11/15 17:17:40 Got Server message [Server First msg] at point (29125, 29250)
2023/11/15 17:17:40 Got Server message [Server Second msg] at point (38083, 38125)
2023/11/15 17:17:40 Got Server message [Server Third msg] at point (41208, 41291)
2023/11/15 17:17:40 Got Server message [Server Fourth msg] at point (43625, 43666)
2023/11/15 17:17:40 Got Server message [Server Fifth msg] at point (45708, 45750)
2023/11/15 17:17:40 Got Server message [Server Sixth msg] at point (47750, 47791)
````