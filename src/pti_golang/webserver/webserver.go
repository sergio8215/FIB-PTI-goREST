package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "encoding/json"
    "io"
    "io/ioutil"
    "encoding/csv"
    "os"
    "bufio"
)

type ResponseMessage struct {
    Field1 string
    Field2 string
}

type RequestMessage struct {
    Field1 string
    Field2 string
}


func main() {
router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/", Index)
router.HandleFunc("/endpoint/{param}", endpointFunc)
router.HandleFunc("/endpoint2/{param}", endpointFunc2JSONInput)

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

func readFiler (w http.ResponseWriter){
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
}




