package cluster

import (
	"glc/cmn"
	"glc/conf"
	"glc/www/service"
	"log"
)

func Start() {
	if !conf.IsClusterMode() {
		return
	}

	log.Println("集群节点启动", cmn.GetLocalGlcUrl())
	joinCluster()
	kv, err := service.GetSysmntItem(KEY_CLUSTER)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(kv.ToJson())
	}

	// 异步检查更新数据
	go dataAsync()
}
