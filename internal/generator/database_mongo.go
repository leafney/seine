package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// addMongoDatabase 添加MongoDB数据库支持
func (g *ComponentGenerator) addMongoDatabase() error {
	// 创建pkg/mongox目录
	mongoxDir := filepath.Join(g.ProjectPath, "pkg/mongox")
	if err := os.MkdirAll(mongoxDir, 0755); err != nil {
		return err
	}

	// 创建mongo.go文件
	mongoFilePath := filepath.Join(mongoxDir, "mongo.go")

	// 检查文件是否已存在
	if _, err := os.Stat(mongoFilePath); err == nil {
		return fmt.Errorf("mongo.go文件已存在")
	}

	// MongoDB模板内容
	mongoTemplate := `/**
 * @Date:        {{.Date}}
 */

package mongox

import (
	"context"

	"{{.ModulePath}}/config"
	"{{.ModulePath}}/pkg/xlogx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSvc struct {
	*mongo.Database
}

func NewMongoSvc(cfg *config.Config, log *xlogx.XLogSvc, stop chan struct{}) *MongoSvc {
	cfgMongo := cfg.MongoDB
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfgMongo.URI))
	if err != nil {
		log.Fatalf("[Mongo] 连接错误 [%v]", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("[Mongo] ping错误 [%v]", err)
	}

	go func() {
		<-stop
		if err := client.Disconnect(ctx); err != nil {
			log.Errorf("[Mongo] 断开连接错误 [%v]", err)
		} else {
			log.Infoln("[Mongo] 退出成功")
		}
	}()

	db := client.Database(cfgMongo.Database)

	log.Infoln("[Mongo] 加载成功")

	return &MongoSvc{db}
}`

	// 创建模板数据
	data := struct {
		Date       string
		ModulePath string
	}{
		Date:       g.Date,
		ModulePath: g.ModulePath,
	}

	// 解析模板
	tmpl, err := template.New("mongo").Parse(mongoTemplate)
	if err != nil {
		return err
	}

	// 创建文件
	file, err := os.Create(mongoFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	// 更新配置文件
	if err := g.updateConfigForDB("mongodb"); err != nil {
		return err
	}

	// 更新wire.go文件
	if err := g.updateWireFileForDB("mongodb"); err != nil {
		return err
	}

	return nil
}
