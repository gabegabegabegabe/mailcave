package mail

import (
	"github.com/tambchop/mailcave/logging"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	msgCollectionName = "message"
)

// MongoMsg is just a Msg with a MongoDB style _id field.
type MongoMsg struct {
	ID bson.ObjectId `bson:"_id"`
	*Msg
}

// MongoArchive provides storage of mail messages using an underlying MongoDB database.
type MongoArchive struct {
	dbAddr  string
	dbName  string
	session *mgo.Session
	logger  *logging.Logger
}

// NewMongoArchive creates a MongoArchive.
func NewMongoArchive(dbAddr string, dbName string, logger *logging.Logger) *MongoArchive {
	return &MongoArchive{
		dbAddr:  dbAddr,
		dbName:  dbName,
		session: nil,
		logger:  logger,
	}
}

// Open causes the archive to open the connection to its MongoDB database.
func (ma *MongoArchive) Open() error {

	if ma.session != nil {
		ma.logger.Printf("MongoDB session is already active")
		return nil
	}

	ma.logger.Printf("connecting to database %s at %s", ma.dbName, ma.dbAddr)

	session, err := mgo.Dial(ma.dbAddr)
	if err != nil {
		return err
	}
	ma.session = session

	ma.logger.Printf("connected to database")

	return nil
}

// Close causes the archive to close its underlying database connection.
func (ma *MongoArchive) Close() {
	ma.session.Close()
	ma.session = nil
}

// ArchiveMessage saves a message to the archive.
func (ma *MongoArchive) ArchiveMessage(msg *Msg) (string, error) {
	session := ma.session.Copy()
	defer session.Close()
	coll := session.DB(ma.dbName).C(msgCollectionName)

	mm := &MongoMsg{
		ID:  bson.NewObjectId(),
		Msg: msg,
	}
	return string(mm.ID), coll.Insert(mm)
}
