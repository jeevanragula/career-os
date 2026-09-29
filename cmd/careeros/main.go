package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"
    "github.com/jeevanragula/career-os/internal/httpapi"
    "github.com/jeevanragula/career-os/internal/store"
)

func main() {
    addr:=os.Getenv("CAREEROS_HTTP_ADDR");if addr==""{addr=":8080"}
    dsn:=os.Getenv("DATABASE_URL")
    var s *store.Store
    if dsn!="" {var err error;ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second);defer cancel();s,err=store.Open(ctx,dsn);if err!=nil{log.Fatal(err)};defer s.Close()}
    server:=&http.Server{Addr:addr,Handler:httpapi.NewHandler(s),ReadHeaderTimeout:10*time.Second}
    log.Printf("CareerOS listening on %s",addr)
    if err:=server.ListenAndServe();err!=nil&&err!=http.ErrServerClosed{log.Fatal(err)}
}
