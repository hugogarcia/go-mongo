package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hugogarcia/go-mongo/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(session *mgo.Session) *UserController {
	return &UserController{session}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
	}

	uid := bson.ObjectIdHex(id)

	user := models.User{}

	err := uc.session.DB("study").C("users").FindId(uid).One(&user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", userJson)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := models.User{}

	json.NewDecoder(r.Body).Decode(&user)

	user.Id = bson.NewObjectId()
	err := uc.session.DB("study").C("users").Insert(user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", userJson)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
	}

	uid := bson.ObjectIdHex(id)

	err := uc.session.DB("study").C("users").RemoveId(uid)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}