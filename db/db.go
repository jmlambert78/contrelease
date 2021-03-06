package db

import "gopkg.in/mgo.v2"

const (
	Host         = "mongodb://localhost,172.17.0.2:27017"
	Database     = "contrel"
	AuthDatabase = "authdb"
	AuthUserName = "root"
	AuthPassword = "pass"
)

var (
	mainSession *mgo.Session
	mainDb      *mgo.Database
)

type MgoDb struct {
	Session *mgo.Session
	Db      *mgo.Database
	Col     *mgo.Collection
}

func init() {

	if mainSession == nil {

		var err error
		mainSession, err = mgo.Dial(Host)

		if err != nil {
			panic(err)
		}

		mainSession.SetMode(mgo.Monotonic, true)
		mainDb = mainSession.DB(Database)

	}

}

func (this *MgoDb) Init() *mgo.Session {

	this.Session = mainSession.Copy()
	this.Db = this.Session.DB(Database)

	return this.Session
}

func (this *MgoDb) C(collection string) *mgo.Collection {
	this.Col = this.Session.DB(Database).C(collection)
	return this.Col
}

func (this *MgoDb) Close() bool {
	defer this.Session.Close()
	return true
}

func (this *MgoDb) DropDb() {
	err := this.Session.DB(Database).DropDatabase()
	if err != nil {
		panic(err)
	}
}

func (this *MgoDb) RemoveAll(collection string) bool {
	this.Session.DB(Database).C(collection).RemoveAll(nil)
	this.Col = this.Session.DB(Database).C(collection)
	return true
}

func (this *MgoDb) Index(collection string, keys []string) bool {
	index := mgo.Index{
		Key:        keys,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := this.Db.C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return true
}

func (this *MgoDb) IsDup(err error) bool {
	if mgo.IsDup(err) {
		return true
	}
	return false
}
func DbMain() {
	// Database Main Conexion
	Db := MgoDb{}
	Db.Init()
	// Index for users
	keys := []string{"name"}
	Db.Index("users", keys)
	// Index for release
	keys = []string{"sname", "svers", "dest"}
	Db.Index("release", keys)
	// Index for roles
	keys = []string{"name", "dc"}
	Db.Index("roles", keys)
	// Index for roles
}
