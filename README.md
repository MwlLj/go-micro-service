# go-micro-service
golang分布式框架

# zookeeper 安装与启动
1. [地址](http://mirrors.hust.edu.cn/apache/zookeeper/)

# 服务注册说明
## import
* s "github.com/MwlLj/go-micro-service/service_discovery_nocache"
* proto "github.com/MwlLj/go-micro-service/common_proto"

## example
	func main() {
		var conns []proto.CConnectProperty
		conns = append(conns, proto.CConnectProperty{ServerHost: "127.0.0.1", ServerPort: 2182, ServiceId: "server_1"})
		sds := s.New(&s.CInitProperty{
			PathPrefix:   "micro-service",
		 	ServerMode:   s.ServerModeZookeeper,
		 	ServerName:   "testserver",
		 	NodeData:     proto.CNodeData{ServerIp: "127.0.0.1", ServerPort: 50000, ServerUniqueCode: "cacd3aa4-4eb8-4bf6-b967-fbcee5377992", Weight: 1},
		 	Conns:        conns,
		 	ConnTimeoutS: 10})
	 	var _ = sds
	 	for {
	 		time.Sleep(100 * time.Millisecond)
		}
	}

## description
* 服务端只需要调用 New 方法, 就可以将服务信息注册到 zookeeper 中, 内部自动实现重连机制
* 对于类似 http 的服务端, 可以不需要填写 ServerUniqueCode 字段
* ServerUniqueCode 字段是针对利用消息队列做成的服务端(这种服务没有ip和port监听, 是通过第三方的 broker 来做为中介者进行消息传递的), 所以这类的服务 必须要填写 ServerUniqueCode 字段, 而 ServerIp / ServerPort 字段非必填
> 这类服务实现分布式机制的一种方式:
	服务端存在一个 ServerUniqueCode, 客户端连接到服务端时, 服务端将自己的 ServerUniqueCode 发送给 客户端, 客户端每次发送消息时, 将这个 uniqueCode 传递给集群中的服务, 每个服务接收到消息, 都判断一下传来的 uniqueCode 是否与自己的 ServerUniqueCode 相等, 只有相等的时候才会接收, 否则丢弃


# 负载均衡说明
## import
* bl "github.com/MwlLj/go-micro-service/loab_blance"
* proto "github.com/MwlLj/go-micro-service/common_proto"

## example
	func main() {
		var conns []proto.CConnectProperty
		conns = append(conns, proto.CConnectProperty{ServerHost: "127.0.0.1", ServerPort: 2182, ServiceId: "server_1"})
		bls, connChan := bl.New(bl.ServerModeZookeeper, &conns, "micro-service", 10)
		// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmRoundRobin)
		// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmWeightRoundRobin)
		// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmRandom)
		// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmWeightRandom)
		// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmIpHash)
		// algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmUrlHash)
		algorithm := bls.GetNormalNodeAlgorithm(bl.AlgorithmLeastConnect)
		select {
		case <-connChan:
			break
		}
		data, err := bls.GetMasterNode("testserver")
		if err != nil {
			fmt.Println("[ERROR] get master node error, ", err)
		} else {
			fmt.Println("------------master------------")
			fmt.Println(data.ServerIp, data.ServerPort, data.ServerUniqueCode, data.Weight)
			fmt.Println("------------master------------")
		}
		fmt.Println("------------normal------------")
		for i := 0; i < 20; i++ {
			// data, err = algorithm.Get("testserver", "192.168.9.2")
			// data, err = algorithm.Get("testserver", "192.168.9.2")
			// data, err = algorithm.Get("testserver", "/data/user")
			data, err = algorithm.Get("testserver", "/data/video")
			if err != nil {
				fmt.Println("[ERROR] get normal node error, ", err)
			} else {
				fmt.Println(data.ServerIp, data.ServerPort, data.ServerUniqueCode, data.Weight)
			}
		}
		fmt.Println("------------normal------------")
		var _ = bls
		for {
			time.Sleep(100 * time.Millisecond)
		}
	}

## description
* ServerModeZookeeper 这种方式是通过接口直接调用的, 后续将支持 Http 版本的
* 通过 GetMasterNode 接口获取 master 节点
> master节点的作用:
	如果需要同步操作数据库等不可以同时进行的事务时, 就不能让多个 server 同时执行, 这种任务, 同一时刻只能有一个服务运行
	这个服务就是 master (master 节点的逻辑不需要使用者关心, 使用者只需要通过 GetMasterNode 接口获取此时的 master 节点对应的服务信息就可以)
* 负载均衡 获取normal节点 是通过 GetNormalNodeAlgorithm 来指定所需的算法
* 通过 GetNormalNodeAlgorithm 返回的对象, 调用 Get 方法就可以, 获取到一个 normal 节点信息

## 结构说明
* mqtt类服务的负载均衡中间件结构
1. 负载均衡中间件服务器
	> 接收所有的消息, 将发送者的url中的最后一部分(微服务的uniquecode) 截取出来, 换成配置文件中匹配的服务的 uuid, 并在后面加上发送者的 topic
	> 流程举例:
		发送者topic:
			/login/user
		接收请求的微服务的名称:
			cfgs
		接收请求的微服务的uniquecode:
			efda4376-529c-4752-a264-553aa0c7f989
		负载均衡中间件服务器接收到的完整topic:
			/login/user/efda4376-529c-4752-a264-553aa0c7f989
		负载均衡中间件服务器发送的topic:
			配置文件匹配的服务在服务中心中注册的服务uniquecode + /login/user
2. 微服务
	> 订阅注册到服务中心的服务uniquecode + topic
3. 目的
