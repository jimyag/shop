package initialize

import (
	"io/ioutil"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

func LocalConfigToRemoteConfig(consulAddress, debugPath, releasePath, filename string) {
	value, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("读取 %s失败 :%s\n", filename, err.Error())
	}
	consulClient, err := consulapi.NewClient(&consulapi.Config{Address: consulAddress})
	if err != nil {
		log.Fatalln("consul连接失败:", err)
	}
	kp := consulapi.KVPair{
		Key:   debugPath,
		Value: value,
	}
	consulClient.KV().Put(&kp, nil)
	kp.Key = releasePath
	consulClient.KV().Put(&kp, nil)
}
