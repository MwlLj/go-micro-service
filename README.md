# go-micro-service
golang分布式框架

# zookeeper 安装与启动
1. [地址](http://mirrors.hust.edu.cn/apache/zookeeper/)

# 服务端说明
## import
* 1 s "github.com/MwlLj/go-micro-service/service_discovery_nocache"
* 2 proto "github.com/MwlLj/go-micro-service/common_proto"
* 3 example
	> var conns []proto.CConnectProperty
	> conns = append(conns, proto.CConnectProperty{ServerHost: "127.0.0.1", ServerPort: 2182, ServiceId: "server_1"})
	> sds := s.New(&s.CInitProperty{
	>	PathPrefix:   "micro-service",
	> 	ServerMode:   s.ServerModeZookeeper,
	> 	ServerName:   "testserver",
	> 	NodeData:     proto.CNodeData{ServerIp: "127.0.0.1", ServerPort: 50000, ServerUniqueCode: "cacd3aa4-4eb8-4bf6-b967-fbcee5377992", Weight: 1},
	> 	Conns:        conns,
	> 	ConnTimeoutS: 10})
	> 	var _ = sds
	> 	for {
	> 		time.Sleep(100 * time.Millisecond)
	>	}
