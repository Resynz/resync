/**
 * @Author: Resynz
 * @Date: 2021/12/31 16:27
 */
package main

import (
	"log"
	"resync/db"
	"resync/queue"
	"resync/server"
)

func main() {
	log.Println("init db handler ...")
	if err:=db.InitDbHandler();err!=nil{
		log.Fatalf("init db failed! err:%s\n",err.Error())
	}
	log.Println("init db tables ...")
	if err:=db.InitDBTables();err!=nil{
		log.Fatalf(err.Error())
	}
	log.Println("init default admin ...")
	if err:=db.InitDefaultAdmin();err!=nil{
		log.Fatalf("init default admin failed! err:%s\n",err.Error())
	}

	log.Println("init queues ...")
	queue.InitQueues()

	server.StartServer()
}
