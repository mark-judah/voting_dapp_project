package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"vote_backend/models"
	"vote_backend/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//var Transactions = make(map[string]models.Transaction)

func AppendToLeader(newVote models.Transaction) {

	transactionData, err3 := json.Marshal(newVote)
	if err3 != nil {
		panic(err3)
	}

	fmt.Println("Writing to db" + fmt.Sprintf("%+v", newVote))
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//commit to log.json file
	PersistLog(newVote)

	//decrease leader log counter to match map size
	//verify the transactions...call  a vote verification function, return bool
	//ensure that the txid is valid, and was generated by an official client app
	//ensure that the node id is valid
	//ensure that the candidate id is valid
	//ensure that the vote hash is valid
	//ensure that the voter exists and that the voters details hash matches the stored hash
	//ensure that the voter hasnt already voted
	//insert verified transaction into db
	database.Create(&newVote)

	token := Client[0].Publish("followerAppend/1", 0, false, transactionData)
	token.Wait()
}

func PersistLog(newVote models.Transaction) {
	data, err := json.MarshalIndent(newVote, "", " ")
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	check_file, err2 := os.Stat("log.json")
	if err2 != nil {
		panic(err2)
	}
	if check_file.Size() == 0 {
		_, err3 := logFile.WriteString("[" + "\n" + string(data) + "\n" + "]")
		if err3 != nil {
			panic(err3)
		}
	} else {
		//delete closing ] first
		b, err3 := os.ReadFile("log.json")
		if err3 != nil {
			panic(err3)
		}
		//convert file to string
		stringJson := string(b)
		//remove ]
		c := strings.Replace(stringJson, "]", "", -1)
		//delete the file
		err4 := os.Remove("log.json")
		if err4 != nil {
			panic(err4)
		}
		//recreate the file
		newFile, err5 := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err5 != nil {
			panic(err5)
		}
		_, err6 := newFile.WriteString(c)
		if err6 != nil {
			panic(err6)
		}
		_, err7 := newFile.WriteString("," + string(data) + "\n" + "]")
		if err7 != nil {
			panic(err7)
		}
	}

}

func NodeSync() {
	//if no log file exists, fetch leader log.json and loop each transaction
	_, err := os.ReadFile("log.json")
	if err != nil {
		//no file
		utils.SetRaftState("syncing")
		token := Client[0].Publish("leaderLogRequest/1", 0, false, "requesting leader log")
		token.Wait()
	}

}
