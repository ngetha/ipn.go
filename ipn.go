package main

import (
    "fmt"
    "net/http"
	"strconv"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	// port
	HttpPort = 8080
	
	// sp Id
	Username = "saf0722"
	
	// sp Password
	Password = "SEndM3MyM0N3y2014"	
	
	// Receive OK
	IPNReceiveOK = "OK|Thanks for sending the $$."
	
	// Receive !OK
	IPNReceiveFail = "FAIL|Grr! Could not receive your $$"	
	
	// Auth Fail
	IPNAuthFail = "Incorrect Username or Password"
	
)

// The IPNMessage type
type IPNRequest struct{
	Id, Orig, Dest, Tstamp, Text, User, Pass, MpesaCode, MpesaAcc, MpesaMsisdn, MpesaTrxDate, MpesaTrxTime, MpesaAmt, MpesaSender string
}

// verify auth status
func authenticate(user, pass string) (auth bool){
	log.Printf("Auth for IPN Req w/ User -> %v", user)
	
	if user == Username && pass == Password {
		auth = true
		log.Println("Auth OK")
	} else {
		auth = false
		log.Println("Auth !OK")
	}
	
	return
}

// a closure with the ipnRequestHandler
func ipnRequestHandler(db *sql.DB) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// grab all the URL params
		user, pwd := r.FormValue("user"), r.FormValue("pass")		
		
		// auth
		if !authenticate(user, pwd) {
			http.Error(w, IPNAuthFail, 403)
			return;
		}
		
		// save to db
		row := db.QueryRow("select proc_gpp from proc_gpp($1);","x")
		exists := false
		if err := row.Scan(&exists); err != nil {
			log.Fatal("Could not connect to DB" , err.Error())
			fmt.Fprintf(w, "y")
			return
		}		
		
		// respond
		return
	})
}


func main(){
	log.Printf("Connecting to DB")
	connectString :=  "user=postgres password=? dbname=? sslmode=disable host=? port=?"
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		log.Fatal("Could not connect to DB" , err)
	}
	log.Println("Got connection to db")	
				
	log.Printf("Setting up HttpHandlers")
	http.Handle("/ipn/mpesa/accept", ipnRequestHandler(db))
	log.Printf("Listening at port %d", HttpPort)
    http.ListenAndServe(":" + strconv.Itoa(HttpPort), nil)
}