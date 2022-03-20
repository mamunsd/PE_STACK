<?PHP
require_once 'phpStack/vendor/autoload.php';
use \MongoDB\BSON\ObjectID as MongoId;

function countRecords() {
    $client = new MongoDB\Client;
    $db = $client->kothok;
    $coll = $db->short_messages;
    $mycursor = $coll->count();
    // $mycursor = $coll->count(['ROOM_NAME'=>'দ্বিতীয় ব্যাচ']);
    return $mycursor;
}

function insertNewSmsg($smsg) {
    $client = new MongoDB\Client;
    $db = $client->kothok;
    $coll = $db->short_messages;
    $myresult = $coll->insertOne($smsg);
    return $myresult;
}

function readSmsg() {
    $client = new MongoDB\Client;
    $db = $client->kothok;
    $coll = $db->short_messages;
    // $mycursor = $coll->find([],  ['limit' => 1000, 'sort' => ['_id' => -1]]);
    $mycursor = $coll->find([],  ['limit' => 2, 'sort' => ['_id' => -1]]);
    // $mycursor = $coll->find([],  ['limit' => 1000, 'sort' => ['_id' => -1]]);
    return json_encode($mycursor->toArray(), JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
}

function readSmsgWp($param) {
    $client = new MongoDB\Client;
    $db = $client->kothok;
    $coll = $db->short_messages;
    $mycursor = $coll->find([],  ['limit' => (int)$param['limit'], 'skip' => (int)$param['skip'], 'sort' => ['TS' => -1]]);
    return $mycursor->toArray();
}

function addAidColl() {
    $client = new MongoDB\Client;
    $db = $client->kothok;
    $coll = $db->short_messages;
    $mycursor = $coll->find([]);
    // return $mycursor->toArray();
    $records = $mycursor->toArray();
    // $total = $coll->count();
    // echo $total . PHP_EOL;
    // var_dump((string)$records[0]['_id']);
    foreach($records as $key=>$value) {
        // $id = (string)$value['_id'];
        $oid = $value['_id'];
        $iid = $key + 1;
        // echo $iid;
        $coll->updateOne(['_id'=>$oid], ['$set'=>['IID'=>$iid]]);
    } 
}



function gnrlgrywithid($qparams) {
    $client = new MongoDB\Client("mongodb://172.16.31.45:27017");   

    // $q_collection = "short_messages";
    $q_collection = $qparams['q_collection'];
    $db = $client->bbloger;
    $coll = $db->short_messages;
    $coll = $db->$q_collection;

    if(array_key_exists("_id", $qparams['filter'])) {
        if(array_key_exists('$gt', $qparams['filter']['_id'])) {
            $id = $qparams['filter']['_id']['$gt'];
            $qparams['filter']['_id']['$gt'] = new MongoId($id);
        }

        if(array_key_exists('$lt', $qparams['filter']['_id'])) {
            $id = $qparams['filter']['_id']['$lt'];
            $qparams['filter']['_id']['$lt'] = new MongoId($id);
        }        
    }


    $mycursor = $coll->find($qparams['filter'], $qparams['qconfig']);
    return $mycursor->toArray();
}

function generalQuery($params) {
    $client = new MongoDB\Client;
    $db = $client->kothok;
    $coll = $db->short_messages;
    $mycursor = $coll->find($params['filter'], $params['options']);
    return $mycursor->toArray();
}

function mongo_CountQuery($query) {
    $client = new MongoDB\Client;
    $db = $client->kothok;
    $coll = $db->short_messages;    
    $mycursor = $coll->find($query['filter'], $query['options']);
    return $mycursor->toArray();
}

function getMsgByRoomLastFifty($params) {
    $client = new MongoDB\Client;
    $db = $client->kothok;
    $coll = $db->short_messages;
    $mycursor = $coll->find(['ROOM_NAME'=>$params['ROOM_NAME']], ['sort'=>['_id'=> -1], 'limit'=>100, 'skip'=>0]);
    return $mycursor->toArray();
}