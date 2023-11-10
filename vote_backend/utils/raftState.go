package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func ReadClientID() string {
	if _, err := os.Stat("clientId"); err == nil {
		fmt.Printf("File exists\n")

		clientId, err := os.ReadFile("clientId")
		if err != nil {
			panic(err)
		}
		return string(clientId)
	} else {
		fmt.Printf("File does not exist\n")

		clientId := []byte(uuid.New().String())
		err := os.WriteFile("clientId", clientId, 0644)
		if err != nil {
			panic(err)
		}
		return string(clientId)
	}

}

func GetClientTerm() int {
	key := ReadClientID() + "term"
	val, err := RedisClient.Get(Ctx, key).Result()
	if err != nil {
		panic(err)
	}
	intVal, err2 := strconv.Atoi(val)
	if err2 != nil {
		panic(err2)
	}
	return intVal
}

func GetClientState() string {
	key := ReadClientID() + "state"
	val, err := RedisClient.Get(Ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func GetClientVote() []string {
	key := ReadClientID() + "votePayload"
	val, err := RedisClient.Get(Ctx, key).Result()
	if err != nil {
		panic(err)
	}
	var dataArray []string
	err = json.Unmarshal([]byte(val), &dataArray)
	if err != nil {
		log.Println("Error unmarshalling:", err)
	}

	return dataArray
}

func SetRaftTerm(term int) {
	err := RedisClient.Set(Ctx, ReadClientID()+"term", term, 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := ReadClientID() + "term"
		val, err := RedisClient.Get(Ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Term", val)
	}
}

func SetRaftState(state string) {
	err := RedisClient.Set(Ctx, ReadClientID()+"state", state, 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := ReadClientID() + "state"
		val, err := RedisClient.Get(Ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("State", val)
	}
}

func SetVoteAndTerm(candidateNodeId string, term string, vote string) {
	fmt.Println("Storing Vote: " + "candidateNodeId: " + candidateNodeId + " term: " + term + " vote: " + vote)
	var votePayload []string
	votePayload = append(votePayload, candidateNodeId, term, vote)
	jsonData, err2 := json.Marshal(votePayload)
	if err2 != nil {
		panic(err2)
	}

	err := RedisClient.Set(Ctx, ReadClientID()+"votePayload", jsonData, 0).Err()
	if err != nil {
		panic(err)
	} else {
		val, err := RedisClient.Get(Ctx, ReadClientID()+"votePayload").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Stored vote", val)
	}
}
