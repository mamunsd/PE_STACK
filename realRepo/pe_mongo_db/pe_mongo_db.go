/*
D: An ordered representation of a BSON document (slice)
M: An unordered representation of a BSON document (map)
A: An ordered representation of a BSON array
E: A single element inside a D type

bson.M হইলো জোসোন এর মতো
bson.D হইলো এরের মতো
bson.E হইলো এক টা জেসোন আইটেম
*/

package pe_mongo_db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mamunsd/PE_STACK/pe_data_models"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var myDbClient *mongo.Client

// var myCollectionPointer *mongo.Collection
var ctx context.Context
var connerr error
var DbName string = "smanager"

/*
Go তে init() ফাংশনটা প্যাকেজ লেভেলে একবার লেখা যায় .. এবং যতবার ই প্যাকেজ টা ইমপোর্ট করা হউক
এই ফাংশন টা এক বার অটোমেটিক কল হবে!!
এই ফাংশনকে আলাদা করে কল করা লাগে  না!!
*/

func init_1() {
	ctx = context.TODO()
	myDbClient, connerr = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if connerr != nil {
		panic(connerr)
	}
	if err := myDbClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
}

func init() {
	ctx = context.TODO()
	myDbClient, connerr = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://172.16.25.101:27017,localhost:28018/?replicaSet=ultraMongoRepl_00&connectTimeoutMS=300000"))
	if connerr != nil {
		panic(connerr)
	}
	if err := myDbClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
}

func PeRandInt(mmin int, mmax int) uint32 {
	rand.Seed(time.Now().UnixNano())
	mrandnum := rand.Intn(mmax-mmin) + mmin
	return uint32(mrandnum)
}

func PeJsonArrayToBsonRecur(myJson string, myPtr *bson.A) {
	mgj := gjson.Parse(myJson)

	mgj.ForEach(func(key, value gjson.Result) bool {
		if !value.IsArray() && !value.IsObject() {
			mydtype := value.Type.String()
			if mydtype == "String" {
				*myPtr = append(*myPtr, value.String())
			} else if mydtype == "Number" {
				*myPtr = append(*myPtr, value.Num)
			}
		} else if value.IsArray() {
			mySarr := bson.A{}
			PeJsonArrayToBsonRecur(value.String(), &mySarr)
			*myPtr = append(*myPtr, mySarr)
		} else if value.IsObject() {
			mySobj := bson.M{}
			PeJsonToBsonRecur(value.String(), &mySobj)
			*myPtr = append(*myPtr, mySobj)
		}
		return true
	})
}

func PeJsonToBsonRecur(myJson string, myPtr *primitive.M) {
	myM := *myPtr
	mgj := gjson.Parse(myJson)
	mgj.ForEach(func(key, value gjson.Result) bool {
		if !value.IsObject() && !value.IsArray() {
			mydtype := value.Type.String()
			if mydtype == "String" {
				thisKey := key.String()
				thisVal := value.String()
				myM[thisKey] = thisVal
			} else if mydtype == "Number" {
				thisKey := key.String()
				thisVal := value.Num
				myM[thisKey] = thisVal
			}

		} else if value.IsArray() {
			thisArray := bson.A{}
			PeJsonArrayToBsonRecur(value.String(), &thisArray)
			myM[key.String()] = thisArray
		} else if value.IsObject() {
			myNO := bson.M{}
			myM[key.String()] = myNO
			PeJsonToBsonRecur(value.String(), &myNO)
		}
		return true
	})
}

func PeJsonFilterToBsonFilter(filterJson string, myPtr *bson.M) {
	myBsonFilter := *myPtr

	myParsed := gjson.Parse(filterJson)
	myParsed.ForEach(func(key, value gjson.Result) bool {
		if key.String() == "_id" {
			myDataType := value.Type.String()
			if myDataType == "String" {
				docID, err := primitive.ObjectIDFromHex(value.String())
				if err != nil {
				}
				myBsonFilter[key.String()] = docID
			} else {
				if value.Get("$gt").Exists() {
					idStr := value.Get("$gt").String()
					docID, err := primitive.ObjectIDFromHex(idStr)
					if err != nil {
					}
					myBsonFilter["_id"] = bson.M{"$gt": docID}
				} else if value.Get("$lt").Exists() {
					idStr := value.Get("$lt").String()
					docID, err := primitive.ObjectIDFromHex(idStr)
					if err != nil {
					}
					myBsonFilter["_id"] = bson.M{"$lt": docID}
				}
			}
		} else if key.String() == "$and" {
			if value.IsArray() {
				myValArray := bson.A{}
				PeJsonArrayToBsonRecur(value.String(), &myValArray)
				myBsonFilter["$and"] = myValArray
			}
		} else if !value.IsObject() && !value.IsArray() {
			myDataType := value.Type.String()
			if myDataType == "String" {
				myBsonFilter[key.String()] = value.String()
			} else if myDataType == "Number" {
				myBsonFilter[key.String()] = value.Num
			}
		}
		return true
	})
}

func UpdateOne(collName string, recordID string, valJson string) {
	mDbCollection := myDbClient.Database(DbName).Collection(collName)
	mFilter := bson.M{}
	recordObjectID, _ := primitive.ObjectIDFromHex(recordID)
	mFilter["_id"] = recordObjectID
	mVal := bson.M{}
	PeJsonToBsonRecur(valJson, &mVal)
	mUpdate := bson.M{}
	mUpdate["$set"] = mVal
	_, err := mDbCollection.UpdateOne(ctx, mFilter, mUpdate)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertOne(collName string, valJson string) {
	mDbCollection := myDbClient.Database(DbName).Collection(collName)

	mVal := bson.M{}

	PeJsonToBsonRecur(valJson, &mVal)
	_, err := mDbCollection.InsertOne(ctx, mVal)

	if err != nil {
		log.Fatal(err)
	}
}

func GenQueryMongo(myQuery []byte) string {
	// const myJson = `{"q_collection":"short_messages","filter":{},"qconfig":{"sort":{"_id":-1},"limit":100}}`
	// const myJson = `{"q_collection":"short_messages","filter":{"_id": {"$gt": "6077cb85fbcb852f9f69ada5"}, "IID" : {"$lt" : 842}, "USER_ID" : "rasel@gmail.com", "CHANNEL" : "PYTHON"},"qconfig":{"sort":{"_id":-1},"limit":75}}`
	// const myJson = `{"q_collection":"short_messages","filter":{"_id": {"$gt": "6077cb85fbcb852f9f69ada5"}, "USER_ID" : "rasel@gmail.com", "CHANNEL" : "PYTHON"},"qconfig":{"sort":{"_id":-1},"limit":40}}`
	// const myJson = `{"q_collection":"short_messages","filter":{"USER_ID":"romanahmed98@gmail.com","ROOM_NAME":"প্রথম ব্যাচ","CHANNEL":"PYTHON","_id":{"$lt":"60832e9afbcb8525ba5b8329"}},"qconfig":{"sort":{"_id":-1},"limit":20}}`
	fmt.Println(string(myQuery))
	bSonFilter := bson.M{}
	m := gjson.Parse(string(myQuery))
	collName := m.Get("q_collection").String()
	filter := m.Get("filter")

	filter.ForEach(func(key, value gjson.Result) bool {
		if !value.IsObject() && !value.IsArray() {
			thisKey := key.String()
			thisVal := value.String()
			bSonFilter[thisKey] = thisVal
		} else {
			// যদি প্রথম লেভেলের পরে ও আবার ভ্যালু আবার জেসোন হয়
			if value.IsObject() {
				// যদি কি টা  মঙ্গো আইডি হয় .. মঙ্গো আইডি এক পদের বাল
				if key.String() == "_id" {
					if value.Get("$gt").Exists() {
						idStr := value.Get("$gt").String()
						docID, err := primitive.ObjectIDFromHex(idStr)
						if err != nil {
						}
						bSonFilter["_id"] = bson.M{"$gt": docID}
					}
					if value.Get("$lt").Exists() {
						idStr := value.Get("$lt").String()
						docID, err := primitive.ObjectIDFromHex(idStr)
						if err != nil {
						}
						bSonFilter["_id"] = bson.M{"$lt": docID}
					}
					// আর যদি দ্বিতীয় লেভেলে Key টা _id মানে বালের মঙ্গো আইডি না হয় .. সেই ক্ষেত্রে $gt এবং $lt অর্থাৎ গ্রেটার আর লেস এর জন্য
				} else {
					if value.Get("$gt").Exists() {
						bSonFilter[key.String()] = bson.M{"$gt": value.Get("$gt").Num}
					}
					if value.Get("$lt").Exists() {
						bSonFilter[key.String()] = bson.M{"$lt": value.Get("$lt").Num}
					}
				}
			} else {
				println("This is not object")
			}
		}
		return true
	})

	usersCollection := myDbClient.Database("kothok").Collection(collName)
	opts := options.Find()
	confopts := m.Get("qconfig")

	confopts.ForEach(func(key, value gjson.Result) bool {

		if !value.IsObject() && !value.IsArray() {
			switch key.String() {
			case "limit":
				opts.SetLimit(value.Int())
			}
		} else {
			if value.IsObject() {
				switch key.String() {
				case "sort":
					sortopts := bson.D{}
					value.ForEach(func(skey, svalue gjson.Result) bool {
						sortopts = append(sortopts, bson.E{skey.String(), svalue.Int()})
						return true
					})
					opts.SetSort(sortopts)
				}
			}
		}
		return true
	})

	cursor, err := usersCollection.Find(ctx, bSonFilter, opts)
	if err != nil {
		panic(err)
	}
	// Go তে এই ভাবে যদি স্লাইস / এরে ইনিশিয়ালাইজ করা হয় তাইলে জেসন মার্শালিং এর সময় নাল রিটার্ন করে
	// var results []KthkMsgs
	// এই জন্য মেক দিয়ে ইনিশিয়ালাইজ করা উচিৎ
	results := make([]pe_data_models.KthkMsgs, 0)

	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	// fmt.Println(results)
	myBytes, err := json.Marshal(results)
	myResult := string(myBytes)
	return string(myResult)
}

func PeMongoGeneralQuery(QueryJson string) (*mongo.Cursor, error) {

	mParsed := gjson.Parse(QueryJson)
	collName := mParsed.Get("q_collection").String()
	jsonFilter := mParsed.Get("filter").String()

	bSonFilter := bson.M{}

	PeJsonFilterToBsonFilter(jsonFilter, &bSonFilter)

	myDbCollection := myDbClient.Database(DbName).Collection(collName)

	/*
		mmFilterJson := `{"$and" : [{"R_RATING" : {"$gt" : 1200}}, {"R_RATING" : {"$lt" : 3000}}]}`
		map[$and:[map[R_RATING:map[$gt:1200]] map[R_RATING:map[$lt:3000]]]]
		[{$and [[{R_RATING [{$gt 1200}]}] [{R_RATING [{$lte 3000}]}]]}]
	*/

	opts := options.Find()

	confopts := mParsed.Get("qconfig")

	confopts.ForEach(func(key, value gjson.Result) bool {

		if !value.IsObject() && !value.IsArray() {
			switch key.String() {
			case "limit":
				opts.SetLimit(value.Int())
			}
		} else {
			if value.IsObject() {
				switch key.String() {
				case "sort":
					sortopts := bson.D{}
					value.ForEach(func(skey, svalue gjson.Result) bool {
						sortopts = append(sortopts, bson.E{skey.String(), svalue.Int()})
						return true
					})
					opts.SetSort(sortopts)
				}
			}
		}
		return true
	})

	myCursor, err := myDbCollection.Find(ctx, bSonFilter, opts)

	if err != nil {
		panic(err)
	}
	return myCursor, err
}

func ShrtMsgGenQ(queryS string) string {
	myCursor, err := PeMongoGeneralQuery(queryS)

	// Go তে এই ভাবে যদি স্লাইস / এরে ইনিশিয়ালাইজ করা হয় তাইলে জেসন মার্শালিং এর সময় নাল রিটার্ন করে
	// var results []KthkMsgs
	// এই জন্য মেক দিয়ে ইনিশিয়ালাইজ করা উচিৎ
	results := make([]pe_data_models.KthkMsgs, 0)

	if err = myCursor.All(ctx, &results); err != nil {
		panic(err)
	}
	myReturn, _ := json.Marshal(results)
	return string(myReturn)
}
