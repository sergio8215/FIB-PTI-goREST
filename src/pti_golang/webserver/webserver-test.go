package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "encoding/json"
   // "io"
   // "io/ioutil"
    "encoding/csv"
    "os"
 //   "bufio"
)

type ResponseMessage struct {
    price float32
    makes string
    model string
    nodias float32
    nounits float32
}

type RequestMessage struct {
    makes string
    model string
    nodias float32
    nounits float32
}


func main() {
router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/", Index)
router.HandleFunc("/endpoint/{param}", endpointFunc)
router.HandleFunc("/endpoint2", endpointFunc2JSONInput)
router.HandleFunc("/rentals", rentalFunc)

log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Service OK")
}

func endpointFunc(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    param := vars["param"]
    res := ResponseMessage{model: "Text1", makes: param}
    json.NewEncoder(w).Encode(res)
}
func endpointFunc2JSONInput(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Service OK")
}

func rentalFunc(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintln(w, "Service OK")
}


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

