package httpapi

import (
    "fmt"
	"net/http"
)

func listPrintersHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    // w.Write([]byte("list\n"))
    fmt.Printf("list\n")
    fmt.Fprintf(w, "list\n")
    return
}

func createPrinterHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("create\n"))
    fmt.Printf("create\n")
    return
}

func showPrinterByIDHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("show by id\n"))
    return
}

func deletePrinterByIDHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("delete by id\n"))
    return
}
