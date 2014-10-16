package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "html/template"
  "regexp"
  "flag"
  "encoding/json"
  "bytes"
)

type Page struct {
    Title string
    Body  []byte
}

type DatapoaResultFields struct {
    Type string
    Id string
}

type DatapoaResultRecords struct {
    Numero string
    Bairro string
    //Endereco string
    Ddd string
    Escola string
    Telefone string
    Cep string
    Id int64
    Email string
    //Localizacao string
}


type DatapoaLink struct {
    Start string
    Next string
}

type DatapoaResult struct {
    Resource_id string
    Fields []*DatapoaResultFields
    Records []*DatapoaResultRecords
    Links *DatapoaLink
    Total int64
}

type DatapoaMessage struct {
    Help string
    Success bool
    Result *DatapoaResult

}

var templates = template.Must(template.ParseFiles("datapoaview.html"))

var validPath = regexp.MustCompile("^/(datapoaview)?/$")

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func datapoaviewHandler(w http.ResponseWriter, r *http.Request) {
    datapoaMessage, _ := getEscolasParticulares()
    err := templates.ExecuteTemplate(w, "datapoaview.html", &datapoaMessage)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        fmt.Print(r.URL.Path)
        if m == nil {
            fmt.Print("Holy shit")
            http.NotFound(w, r)
            return
        }
        fmt.Print("Hi crasy folks")
        fn(w, r)
    }
}

func main() {
    flag.Parse()
    http.HandleFunc("/datapoaview/", makeHandler(datapoaviewHandler))
    http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func getEscolasParticulares() (DatapoaMessage, error) {
  url := "http://datapoa.com.br/api/action/datastore_search?resource_id=8f8b8f0a-45eb-4372-94e3-058723082a28"
  resp, err := http.Get(url)
  if err != nil {
    fmt.Print("Shit")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  var m DatapoaMessage;
  err = json.Unmarshal(formatEntidates([]byte(body)), &m)
  return m,  err
}

func formatEntidates(b []byte) ([]byte) {
  b = bytes.Replace(b, []byte("\"N\u00daMERO\""),[]byte("\"numero\""), -1)
  b = bytes.Replace(b, []byte("\"BAIRRO\""),[]byte("\"bairro\""),  -1)
  //b = bytes.Replace(b, []byte("\"ENDERE\u00c7O\""),[]byte("\"endereco\""),-1)
  b = bytes.Replace(b, []byte("\"DDD\""),[]byte("\"ddd\""), -1)
  b = bytes.Replace(b, []byte("\"ESCOLA\""),[]byte("\"escola\""), -1)
  b = bytes.Replace(b, []byte("\"TELEFONE\""),[]byte("\"telefone\""), -1)
  b = bytes.Replace(b, []byte("\"CEP\""),[]byte("\"cep\""),-1)
  b = bytes.Replace(b, []byte("\"_id\""),[]byte("\"id\""), -1)
  b = bytes.Replace(b, []byte("\"EMAIL\""),[]byte("\"email\""), -1)
  //b = bytes.Replace(b, []byte("\"LOCALIZA\u00c7\u00c3O\""),[]byte("\"localizacao\""), -1)
  //b = bytes.Replace(b, []byte("\"_links\""),[]byte("\"links\""), -1)
  fmt.Printf("%s", b)
  return b
}

