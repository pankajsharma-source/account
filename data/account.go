package data

import (
	"encoding/json"
	"errors"
	"github.com/couchbase/gocb/v2"
	"github.com/go-playground/validator/v10"
	"io"
	"time"
)

var (
	cbConnStr  		= "couchbase://localhost"
	cbUsername 		= "Administrator"
	cbPassword 		= "password"
	cbBucket   		= "aiFitnessBucket"
	cbScope	   		= "user"
	cbCollection	= "profile"
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user does not exist")
	ErrBadPassword   = errors.New("password does not match")
	ErrBadAuthHeader = errors.New("bad authentication header format")
	ErrBadAuth       = errors.New("invalid auth token")
)

var globalCluster *gocb.Cluster = nil
var globalBucket *gocb.Bucket = nil
var scope *gocb.Scope = nil
var collection *gocb.Collection = nil

// Product defines the structure for an API product
type User struct {
	ID          		string  `json:"id" validate:"required"`
	FirstName        	string  `json:"first-name" validate:"required"`
	LastName        	string  `json:"last-name" validate:"required"`
	Height				int  	`json:"height" validate:"gt=0"`
	Weight				int  	`json:"weight" validate:"gt=0"`
	Gender       		string  `json:"gender" validate:"required"`
	DOB         		string  `json:"dob" validate:"required,datetime=2006-01-02"`
	Country   			string  `json:"country" validate:"required"`
	Zip   				int  	`json:"zip" validate:"required"`
	Email   			string  `json:"email" validate:"required,email"`
}

func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func CreateConnection() {

	var err error
	// Connect to Couchbase
	// Get a bucket reference
	opts := gocb.ClusterOptions{
		Username: cbUsername,
		Password: cbPassword,
	}

	globalCluster, err = gocb.Connect(cbConnStr, opts)
	if err != nil {
		panic(err)
	}

	// get a bucket reference
	globalBucket = globalCluster.Bucket(cbBucket)

	// We wait until the bucket is definitely connected and setup.
	err = globalBucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	// for userection and scope
	scope = globalBucket.Scope(cbScope)
	collection = scope.Collection(cbCollection)
}


func AddUser(u *User) {

	if globalCluster == nil {
		CreateConnection()
	}
	collection.Upsert(u.ID, u, &gocb.UpsertOptions{})
}

func GetUser(id string) (User, error){

	if globalCluster == nil {
		CreateConnection()
	}
	getResult, err := collection.Get(id, nil)
	if err != nil {
		return User{}, err
	}
	var user User
	err = getResult.Content(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}