package main

import (
	"flag"
	"log"
	"photo-web/bootstrap"
	"photo-web/config"
)

func main() {
	// all 迁移+种子数据。 schema仅迁移表结构。 seed 仅插入种子数据
	var action string
	flag.StringVar(&action, "action", "all", "all | schema | seed")

	var configFile string
	flag.StringVar(&configFile, "config", "config.yaml", "config file path")

	flag.Parse()

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatal("配置加载错误，err：", err)
	}

	db, err := bootstrap.InitDB(cfg)
	if err != nil {
		log.Fatal("数据库连接错误, err:", err)
	}

	switch action {
	case "all":
		err := Migrate(db)
		if err != nil {
			log.Fatal(err)
		}
		err = SeedStyles(db)
		if err != nil {
			log.Fatal(err)
		}
	case "schema":
		err := Migrate(db)
		if err != nil {
			log.Fatal(err)
		}
	case "seed":
		err = SeedStyles(db)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("未知命令：", action)
	}

}
