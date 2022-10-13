// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/segmentio/kafka-go"
// )

// type TopicConfig struct {
// 	Topic        string `json:"topic"`
// 	Partitions   int    `json:"partitions"`
// 	Replications int    `json:"replications"`
// }
// type Topic struct {
// 	Topic string `json:"topic"`
// }

// func producerHandler() func(http.ResponseWriter, *http.Request) {
// 	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
// 		var input Topic
// 		jsonError := json.NewDecoder(req.Body).Decode(&input)
// 		if jsonError != nil {
// 			fmt.Println("ERROR IN BODY")
// 			return
// 		}
// 		fmt.Println(input.Topic)
// 		w := &kafka.Writer{
// 			Addr:  kafka.TCP("kafka:9092"),
// 			Topic: input.Topic,
// 		}

// 		err := w.WriteMessages(context.Background(), kafka.Message{
// 			Key:   []byte("hello"),
// 			Value: []byte("hello from producer !"),
// 		})

// 		if err != nil {
// 			fmt.Println(err)
// 		} else {
// 			fmt.Println("MESSAGE SENT BY PRODUCER!")
// 		}
// 	})
// }
// func createTopic() func(http.ResponseWriter, *http.Request) {
// 	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {
// 		conn, err := kafka.Dial("tcp", "kafka:9092")
// 		var input TopicConfig
// 		jsonError := json.NewDecoder(req.Body).Decode(&input)
// 		if jsonError != nil {
// 			fmt.Println("ERROR IN BODY")
// 			return
// 		}

// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		defer conn.Close()

// 		controller, err := conn.Controller()
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		fmt.Println(controller.Host + ":" + strconv.Itoa(controller.Port))
// 		controllerConn, err := kafka.Dial("tcp", controller.Host+":"+strconv.Itoa(controller.Port))
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		defer controllerConn.Close()

// 		topicConfigs := []kafka.TopicConfig{{Topic: input.Topic, NumPartitions: input.Partitions, ReplicationFactor: input.Replications}}

//			err = controllerConn.CreateTopics(topicConfigs...)
//			if err != nil {
//				fmt.Println(err.Error())
//			} else {
//				fmt.Println("DONE")
//			}
//		})
//	}
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/whitedevil31/atlan-backend/producer-api/routes"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterStudentRoutes(r)
	http.Handle("/", r)
	host := "localhost:8080"
	fmt.Println("APP RUNNING ON " + host)
	http.ListenAndServe(host, r)

}
