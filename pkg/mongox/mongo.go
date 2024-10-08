/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-10-08 18:47
 * @Description:
 */

package mongox

import (
	"context"
	"github.com/leafney/seine/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"
)

type MongoDBSvc struct {
	*mongo.Client
}

func NewMongoDBSvc(cfg *config.Config, stop chan struct{}) *MongoDBSvc {
	cfgMongo := cfg.Mongo

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfgMongo.Addr))
	if err != nil {
		log.Fatalf("[Mongo] connect error [%v]", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("[Mongo] ping error [%v]", err)
	}

	go func() {
		<-stop
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("[Mongo] disconnect error [%v]", err)
		} else {
			log.Println("[Mongo] Exit successful")
		}
	}()

	client.Database(cfgMongo.DB)

	log.Println("[Mongo] Load successful")
	return &MongoDBSvc{client}
}
