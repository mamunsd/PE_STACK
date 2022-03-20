package main

import (
	"fmt"
	"pemgo/peMongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func RecursiveJsonToBsonTest() {
	myBson := bson.M{}
	myJson := `{
		"Name" : "Prosenjit",
		"Father's Name" : "Horidash Pal",
		"Age" : 27,
		"Subjects" : ["Bangla", "English","Gonit", {"sub1" : "Vugol", "sub2" : "Orthoniti"}],
		"Marks" : {"Bangla" : 45, "English" : 60, "Gonit" : 85},
		"FullRec" : {
			"Name" : "Prosenjit",
			"Father's Name" : "Horidash Pal",
			"Age" : 27,
			"Subjects" : ["Bangla", "English","Gonit", {"sub1" : "Vugol", "sub2" : "Orthoniti"}],
			"Marks" : {"Bangla" : 45, "English" : 60, "Gonit" : 85},
			"FullSecRec" : {
				"Name" : "Prosenjit",
				"Father's Name" : "Horidash Pal",
				"Age" : 27,
				"Subjects" : ["Bangla", "English","Gonit", {"sub1" : "Vugol", "sub2" : "Orthoniti"}],
				"Marks" : {"Bangla" : 45, "English" : 60, "Gonit" : 85}
			}
		}
	}`

	myStartTime := time.Now().Unix()

	var mindex uint32
	var myTurn uint32
	myTurn = 10000000
	for mindex = 0; mindex < myTurn; mindex++ {
		peMongo.PeJsonToBsonRecur(myJson, &myBson)
	}

	myEndTime := time.Now().Unix()
	required := (myEndTime - myStartTime)
	perSec := int64(myTurn) / required
	fmt.Println(required, perSec)
	// fmt.Println(myBson)
}

func TestRecJsonToBsonSingle() {
	myJson := `{
		"Name" : "Prosenjit",
		"Father's Name" : "Horidash Pal",
		"Age" : 27,
		"Subjects" : ["Bangla", "English","Gonit", {"sub1" : "Vugol", "sub2" : "Orthoniti"}],
		"Marks" : {"Bangla" : 45, "English" : 60, "Gonit" : 85},
		"FullRec" : {
			"Name" : "Prosenjit",
			"Father's Name" : "Horidash Pal",
			"Age" : 27,
			"Subjects" : ["Bangla", "English","Gonit", {"sub1" : "Vugol", "sub2" : "Orthoniti"}],
			"Marks" : {"Bangla" : 45, "English" : 60, "Gonit" : 85},
			"FullSecRec" : {
				"Name" : "Prosenjit",
				"Father's Name" : "Horidash Pal",
				"Age" : 27,
				"Subjects" : ["Bangla", "English","Gonit", {"sub1" : "Vugol", "sub2" : "Orthoniti"}],
				"Marks" : {"Bangla" : 45, "English" : 60, "Gonit" : 85}
			}
		}
	}`
	myBson := bson.M{}
	peMongo.PeJsonToBsonRecur(myJson, &myBson)
	fmt.Println(myBson)
}

func TestUpdateOne() {
	myJson := `{
		"CONTENT" : "হৃদয়ের গাঁথা মালা গাঁথা হলো না!!!",
		"MONTH" : "SEP",
		"OAKASH" : [
			"PRODIP", "JELO", "NA",
			{
				"SON" : 71,
				"MOVIE" : "কথা ছিলো দেখা হবে, দেখা হলো না"
			}
		]
		}`
	recID := "6077c87cfbcb852f9f69aca2"
	peMongo.UpdateOne("short_messages", recID, myJson)
}

func TestInsertOne() {
	myJson := `{"CONT_NAME" : "Sooryavanshi", "CONT_TYPE" : "MOVIE", "INDUSTRY" : "Bollywood"}`
	peMongo.InsertOne("contents", myJson)
}

func TestGenQuery() {
	// QueryJson := `{"q_collection":"short_messages","filter":{"_id": {"$gt": "6077cb85fbcb852f9f69ada5"}, "USER_ID" : "rasel@gmail.com", "CHANNEL" : "PYTHON"},"qconfig":{"sort":{"_id":-1},"limit":1}}`

	QueryJson := `{"q_collection":"short_messages","filter":{"_id": {"$gt": "6077cb85fbcb852f9f69ada5"}, "CHANNEL" : "PYTHON", "$and" : [{"R_RATING" : {"$gt" : 1200}}, {"R_RATING" : {"$lt" : 3000}}]},"qconfig":{"sort":{"_id":-1},"limit":1}}`
	// QueryJson := `{"_id": {"$gt": "6077cb85fbcb852f9f69ada5"}, "$and" : [{"R_RATING" : {"$gt" : 1200}}, {"R_RATING" : {"$lt" : 3000}}], "USER_ID" : "rasel@gmail.com", "CHANNEL" : "PYTHON"}`

	myRes := peMongo.ShrtMsgGenQ(QueryJson)
	fmt.Println(myRes)
}
