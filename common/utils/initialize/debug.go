package initialize

import (
	"io/ioutil"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

//
// LocalConfigToRemoteConfig
//  @Description: 配置文件上传到consul中
//  @param consulAddress
//  @param debugPath
//  @param releasePath
//  @param filename
//
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
	_, err = consulClient.KV().Put(&kp, nil)
	if err != nil {
		log.Fatalln("上传 debug 配置文件失败", err)
	}
	kp.Key = releasePath
	_, err = consulClient.KV().Put(&kp, nil)
	if err != nil {
		log.Fatalln("上传 release 配置文件失败", err)
	}

}
