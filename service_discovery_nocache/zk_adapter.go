package service_discovery_nocache

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _ = fmt.Println

type CZkAdapter struct {
	m_hostPro *zk.DNSHostProvider
	m_hostMap sync.Map
}

func (this *CZkAdapter) Init(property *CInitProperty) error {
	hosts := this.connProperty2Hosts(&property.ConnProperty)
	this.m_hostPro = new(zk.DNSHostProvider)
	err := this.m_hostPro.Init(*hosts)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// server, retryStart := hostPro.Next() //获得host
	// hostPro.Connected()                  //连接成功后会调用
	return nil
}

func (this *CZkAdapter) Connect() error {
	host, _ := this.m_hostPro.Next()
	v, ok := this.m_hostMap.Load(host)
	if !ok {
		return errors.New("[Error] no vaild host")
	}
	net := v.(CNet)
	conn, connChan, err := zk.Connect([]string{host}, time.Duration(net.ConnTimeoutS))
	if err != nil {
		fmt.Println("connect zookeeper server error")
		return err
	}
	select {
	case event := <-connChan:
		if event.State == zk.StateConnected {
			return nil
		}
	case <-time.After(time.Second * 3):
		return errors.New("[Error] connect timeout")
	}
	_ = conn
	return nil
}

func (this *CZkAdapter) connProperty2Hosts(property *CConnectProperty) *[]string {
	var hosts []string
	for _, net := range property.Nets {
		host := this.joinHost(net.ServerHost, net.ServerPort)
		hosts = append(hosts, host)
		this.m_hostMap.Store(host, net)
	}
	return &hosts
}

func (*CZkAdapter) joinHost(ip string, port int) string {
	return strings.Join([]string{ip, strconv.FormatInt(int64(port), 10)}, ":")
}
