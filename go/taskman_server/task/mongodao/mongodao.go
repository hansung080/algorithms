// Package mongodao implements task.Accessor interface with MongoDB.
package mongodao

import (
	"github.com/hansung080/taskman_server/task"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoAccessor is an Accessor for MongoDB.
type MongoAccessor struct {
	session    *mgo.Session
	collection *mgo.Collection
}

// New returns a new MongoAccessor.
func New(path, db, c string) task.Accessor {
	session, err := mgo.Dial(path)
	if err != nil {
		return nil
	}

	collection := session.DB(db).C(c)
	return &MongoAccessor{
		session:    session,
		collection: collection,
	}
}

func (m *MongoAccessor) Close() error {
	m.session.Close()
	return nil
}

// idToObjectId returns bson.ObjectId converted from task.ID.
func idToObjectId(id task.ID) bson.ObjectId {
	return bson.ObjectIdHex(string(id))
}

// objectIdToID returns task.ID converted from bson.ObjectId.
func objectIdToID(objId bson.ObjectId) task.ID {
	return task.ID(objId.Hex())
}

// Get returns a task with a given ID.
func (m *MongoAccessor) Get(id task.ID) (task.Task, error) {
	t := task.Task{}
	err := m.collection.FindId(idToObjectId(id)).One(&t)
	return t, err
}

// Post adds a new task.
func (m *MongoAccessor) Post(t task.Task) (task.ID, error) {
	objId := bson.NewObjectId()
	_, err := m.collection.UpsertId(objId, &t)
	return objectIdToID(objId), err
}

// Put updates a task with with a given ID with t.
func (m *MongoAccessor) Put(id task.ID, t task.Task) error {
	return m.collection.UpdateId(idToObjectId(id), t)
}

// Delete removes a task with a given ID.
func (m *MongoAccessor) Delete(id task.ID) error {
	return m.collection.RemoveId(idToObjectId(id))
}
