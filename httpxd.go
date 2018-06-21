package main

import (
   "log"
   "fmt"
   "net/http"
)


func errorResponse(writer http.ResponseWriter, code int) {
   const template string = `{"success":0,code:7%d,"message":"%s"}`
   var message string = "uncnown code";
   switch (code) {
   case 404:
      message = "document not found"
      break
   }
   http.Error(writer, fmt.Sprintf(template, code, message), code)
}

func main() {

   http.HandleFunc("/t", func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprint(w, "hello")
   })

   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      if r.URL.Path == "/" {
         fmt.Fprint(w, `{"succes":1,"message":"ok"}`)
      } else {
         errorResponse(w, 404)
      }
   })

   log.Println("Starting server on :8032")
   log.Fatalf("cannot start server %v",
         http.ListenAndServe(":8032", nil))
}