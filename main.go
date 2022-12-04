package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hugogarcia/go-mongo/controllers"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"crypto/tls"
	"net"
)

func main(){
	godotenv.Load()
	router := httprouter.New()	
	uc := controllers.NewUserController(getSession())
	router.GET("/user/:id", uc.GetUser)
	router.POST("/user", uc.CreateUser)
	router.DELETE("/user/:id", uc.DeleteUser)

	fmt.Println("RUNNING ON localhost:9090")
	http.ListenAndServe("localhost:9000", router)	
}

func getSession() *mgo.Session{	
	conn := os.Getenv("DB_STRING")
	fmt.Println(conn)

	dialInfo := &mgo.DialInfo{
		Addrs: []string{os.Getenv("DB_STRING")},
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
		Source: "admin",
		
	}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
        conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
        return conn, err
    }

	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil{
		panic(err)
	}

	return s
}