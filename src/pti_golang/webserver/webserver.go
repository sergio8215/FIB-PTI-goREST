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
    "strconv"
    "bufio"
)

type ResponseMessage struct {
    Price int
    Makes string
    Model string
    Nodias int
    Nounits int
}

type RequestMessage struct {
    Makes string
    Model string
    Nodias int
    Nounits int
}


func main() {
router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/", Index)
router.HandleFunc("/endpoint/{param}", endpointFunc)
router.HandleFunc("/endpoint2/{param}", endpointFunc2JSONInput)
router.HandleFunc("/newrental", rentalFunc)
router.HandleFunc("/listrental", listRentalsFunc)

log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Service OK")
}

func endpointFunc(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    param := vars["param"]
    res := ResponseMessage{Model: "Text1", Makes: param}
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
        fmt.Fprintln(w, "Successfully received request with Field1 =", requestMessage.Model)
        fmt.Println(r.FormValue("queryparam1"))
    }
}

func rentalFunc(w http.ResponseWriter, r *http.Request) {
    var requestMessage RequestMessage
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        fmt.Fprintln(w, "error 1")
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        fmt.Fprintln(w, "error 2")
        panic(err)
    }
    if err := json.Unmarshal(body, &requestMessage); err != nil {
        fmt.Fprintln(w, "error 3")
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    } else {
        //fmt.Fprintln(w, "Service OK")
        price := (requestMessage.Nodias)*(requestMessage.Nounits)*3
        fmt.Fprintln(w, "Price of rental= ", price)
        i := fmt.Sprintf("\n\n Make: %#v \n Model: %#v \n Número de días: %#v \n Número de unidades: %#v \n",requestMessage.Makes,requestMessage.Model,requestMessage.Nodias,requestMessage.Nounits)
        fmt.Println(i)

        file, err := os.OpenFile("rentals.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        if err!=nil {
            json.NewEncoder(w).Encode(err)
            return
        }
        writer := csv.NewWriter(file)
        ndias := strconv.Itoa(requestMessage.Nodias)
        nunits := strconv.Itoa(requestMessage.Nounits)
        var data1 = []string{requestMessage.Makes,requestMessage.Model,ndias,nunits}
        writer.Write(data1)
        writer.Flush()
        file.Close()

    }
}

func listRentalsFunc(w http.ResponseWriter, r *http.Request){
    file, err := os.Open("rentals.csv")    
    if err!=nil {
        json.NewEncoder(w).Encode(err)
        return
    }
    reader := csv.NewReader(bufio.NewReader(file))
    i := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        i++

        fmt.Fprintf(w, "The %d rental is: %q \n", i,record)
    }
}

