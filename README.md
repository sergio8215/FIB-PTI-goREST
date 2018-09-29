# Making a Web API in Go

## 1. Introduction

A web service is a generic term for a software function that is accessible through HTTP. Traditional web services usually relied in support protocols for data exchange (e.g. SOAP) and service definition (WSDL). However, nowadays the paradigm has evolved to a simplified form, usually called web APIs. Web APIs normally rely only in plain HTTP (plus JSON for serializing the messages). Their design is usually influenced by the [REST architectural style](https://en.wikipedia.org/wiki/Representational_state_transfer), though the most part of existing web APIs do not really comply with REST principles. Nowadays, the most part of client-server systems (e.g. web applications and mobile apps) design their back end as a combination of web APIs.  

The goal of this session is to create simple web API with the Go programming language and JSON. We will not bother to follow the REST principles, so it will not be a trully RESTful API.  


## 2. Setup

### 2.1 Booting the machine 

Conventional room: Select a 64-bit Linux image and login with your credentials.

Operating Systems room: Select the latest Ubuntu imatge (e.g. Ubuntu 14) with credentials user=alumne and pwd:=sistemes

### 2.2 Prerequisites

Let's start by updating the Ubuntu's Package Index:

    sudo apt-get update

You need 'curl' installed (check it typing 'curl -V' in the terminal). If, for some strange reason, it's not, just do:

    sudo apt-get install curl

You need 'git' installed (in order to be able to run 'go get' commands):

    sudo apt-get install git

(OPTIONAL) NOTE: If you encounter lock errors with apt-get commands try killing the blocking process:

	ps aux | grep apt
	sudo kill PROCNUM

### 2.3 Install Go

Download Go (if the link does not work go [here](https://golang.org/dl/))
    
    wget https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz

(IMPORTANT: The executable will not work if you are working on a 32-bit Linux!)

Extract it into /usr/local 

    sudo tar -C /usr/local -xzf go1.7.1.linux-amd64.tar.gz

Change permissions:
    
    sudo chmod -R 777 /usr/local/go

Add /usr/local/go/bin to the PATH environment variable:

    export PATH=$PATH:/usr/local/go/bin

### 2.4 Setup a directory hierarchy for your projects

(check [this](https://golang.org/doc/code.html) for more info in how to write Go code)

Create a directory to contain your golang workspace (e.g. $HOME/go): 

    cd
    mkdir $HOME/go

Set the GOPATH environment variable to point to that location

    export GOPATH=$HOME/go

Create a directory to contain source files: 

    mkdir $HOME/go/src

Now we will create a package directory (pti_golang/hello) where we will place our first program:

    mkdir -p $HOME/go/src/pti_golang/hello
    cd $HOME/go/src/pti_golang/hello

Now we download a source file:

    wget https://gitlab.fib.upc.edu/pti/pti/raw/master/p3_goREST/src/hello/hello.go

Then we run the "go install" command. This command (1) builds the package pti_golang/hello, producing an executable binary (using the last name of the package path as name of the command), and (2) copies it to the workspace's bin directory (it is created if it does not exists). 

    go install pti_golang/hello

Now we can execute the binary:

    $HOME/go/bin/hello
  
## 3 A simple web server
    
A web API is a specific type of web (HTTP-based) service. Let's start by programming a basic web server with Go:   

Create a directory for this program:

    mkdir $HOME/go/src/pti_golang/webserver

Edit $HOME/go/src/pti_golang/webserver/webserver.go

    package main

    import (
        "fmt"
        "html"
        "log"
        "net/http"
    )

    func main() {
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
        })

        log.Fatal(http.ListenAndServe(":8080", nil))

    } 

Build (will create an executable within $HOME/go/bin/webserver):

    go install pti_golang/webserver

Run:

    $HOME/go/bin/webserver &

test in browser: http://localhost:8080

Kill the process before going through the next steps.
    
## 4 URL routing
    
An web API exposes different functionalities. These functionalities are accessed through different URL routes or endpoints. We need a mechanism that let us map URL routes into calls to different functions in our code. The standard golang library offers a [too complex routing mechanism](https://husobee.github.io/golang/url-router/2015/06/15/why-do-all-golang-url-routers-suck.html), so we will use an external library for that (mux router from the Gorilla Web Toolkit):

    go get "github.com/gorilla/mux"

(check that a new package object has been created within $HOME/go/pkg).

Let's modify our webserver.go to add some routes:

    package main

    import (
        "fmt"
        "log"
        "net/http"
        "github.com/gorilla/mux"
    )

    func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
    router.HandleFunc("/endpoint/{param}", endpointFunc)

    log.Fatal(http.ListenAndServe(":8080", router))
    }

    func Index(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Service OK")
    }

    func endpointFunc(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        param := vars["param"]
        fmt.Fprintln(w, "You are calling with param:", param)
    }

Rebuild, run and open http://localhost:8080/endpoint/1234 in your browser.
   
## 5. JSON 

Typically an endpoint has to deal with more complex input and output parameters. This is usually solved by formatting the parameters (input and/or output) with JSON. Let's modify our webserver.go to include a JSON response.

    package main

    import (
        "fmt"
        "log"
        "net/http"
        "github.com/gorilla/mux"
        "encoding/json"
    )

    type ResponseMessage struct {
        Field1 string
        Field2 string
    }

    func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
    router.HandleFunc("/endpoint/{param}", endpointFunc)

    log.Fatal(http.ListenAndServe(":8080", router))
    }

    func Index(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Service OK")
    }

    func endpointFunc(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        param := vars["param"]
        res := ResponseMessage{Field1: "Text1", Field2: param}
        json.NewEncoder(w).Encode(res)
    }

Rebuild, run and open http://localhost:8080/endpoint/1234 in your browser.

Let's now add a new endpoint that accepts JSON as input. First of all add the following struct:

    type RequestMessage struct {
        Field1 string
        Field2 string
    }

And "io" and "io/ioutil" to the imports:

	import (
        	"fmt"
        	"log"
        	"net/http"
        	"github.com/gorilla/mux"
        	"encoding/json"
		"io"
		"io/ioutil"
    	) 

Then add a new route:

    router.HandleFunc("/endpoint2/{param}", endpointFunc2JSONInput)

And its related code:

    func endpointFunc2JSONInput(w http.ResponseWriter, r *http.Request) {
        var requestMessage RequestMessage
        body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
        if err != nil {
            panic(err)
        }
        if err := r.Body.Close(); err != nil {
            panic(err)
        }
        if err := json.Unmarshal(body, &requestMessage); err != nil {
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            w.WriteHeader(422) // unprocessable entity
            if err := json.NewEncoder(w).Encode(err); err != nil {
                panic(err)
            }
        } else {
            fmt.Fprintln(w, "Successfully received request with Field1 =", requestMessage.Field1)
            fmt.Println(r.FormValue("queryparam1"))
        }
    }

Rebuild and run. In order to submit a JSON request we will use curl instead of the browser. Open a new terminal and type:

    curl -H "Content-Type: application/json" -d '{"Field1":"Value1", "Field2":"Value2"}' http://localhost:8080/endpoint2/1234?queryparam1=Value3

(while curl is enough for this session, for your project you could take a look at [POSTMAN](https://www.getpostman.com/))

**WARNING: The fields of the request and response structs MUST START WITH A CAPITAL LETTER.**

**WARNING: The fields of the request and response structs MUST ONLY INCLUDE ALPHANUMERIC CHARACTERS (AVOID UNDERSCORES, ETC.).**

   
## 6. Creating your own car rental web API

As an example web API you will create a simple car rental web API. It will consist in two functionalities:

- Request a new rental: An endpoint to register a new rental order. Input fields will include the car maker, car model, number of days and number of units. The total price of the rental will be returned to the user along with the data of the requested rental.
 
- Request the list of all rentals: An endpoint that will return the list of all saved rental orders (in JSON format). 

In order to keep the rentals data (to be able to list them) you will need to save the data to the disk. A single text file where each line represents a rental will be enough (though not in a real scenario). 

## ANNEX 1. Writing comma-separated values to a CSV file

An easy way to save the list of rentals could be a text file with lines containing comma-separated values (CSV). One rental per line. This way you can save a new rental just adding a line at the end of the file.

	(need to add "encoding/csv" and "os" to imports)

	func writeToFile(w http.ResponseWriter) {
		file, err := os.OpenFile("rentals.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err!=nil {
			json.NewEncoder(w).Encode(err)
			return
   		}
		writer := csv.NewWriter(file)
		var data1 = []string{"Toyota", "Celica"}
		writer.Write(data1)
		writer.Flush()
		file.Close()
	}

## ANNEX 2. Reading a CSV file

In order to read the list of rentals from the CSV file you can do:

	(need to add "bufio" to imports)

	file, err := os.Open("rentals.csv")
        if err!=nil {
		json.NewEncoder(w).Encode(err)
		return
    	}
    	reader := csv.NewReader(bufio.NewReader(file))
    	for {
        	record, err := reader.Read()
        	if err == io.EOF {
            		break
                }
                fmt.Fprintf(w, "The first value is %q", record[0])
        }
	
     
## ANNEX 3. Dynamic growing arrays in Go: Slices

If you need reading the rentals file into an array, whose initial size is unknown, you can use a  [Go slice](https://blog.golang.org/go-slices-usage-and-internals) and the "append" function this way: 

    var lines []ResponseMessage
    for {
        ...
	lines = append(lines, ResponseMessage{Field1: "Text1", Field2: "Text2"})
        
    }
    ...


