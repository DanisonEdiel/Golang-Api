package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net"
    "net/http"
    "net/rpc"
)

type Calculator struct{}

func (c *Calculator) Multiply(args *Args, result *int) error {
    if args == nil {
        return errors.New("invalid arguments")
    }
    *result = args.A * args.B
    return nil
}

type Args struct {
    A, B int
}

// RPC Server
func startRPCServer() {
    calculator := new(Calculator)
    err := rpc.Register(calculator)
    if err != nil {
        fmt.Println("Error registering service:", err)
        return
    }

    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        fmt.Println("Error setting up listener:", err)
        return
    }
    defer listener.Close()
    fmt.Println("RPC Server listening on port 1234")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Connection error:", err)
            continue
        }
        go rpc.ServeConn(conn)
    }
}

// HTTP Handler for RPC
func rpcHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    client, err := rpc.Dial("tcp", "localhost:1234")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error connecting to RPC server: %v", err), http.StatusInternalServerError)
        return
    }
    defer client.Close()

    var args Args
    if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
        http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
        return
    }

    var result int
    err = client.Call("Calculator.Multiply", &args, &result)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error calling RPC method: %v", err), http.StatusInternalServerError)
        return
    }

    response := map[string]int{"result": result}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Enable CORS for HTTP Server
func enableCORS(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        h(w, r)
    }
}

// HTTP Server
func startHTTPServer() {
    http.HandleFunc("/rpc/multiply", enableCORS(rpcHandler))
    fmt.Println("HTTP Server listening on port 8081")
    http.ListenAndServe(":8081", nil)
}

func main() {
    go startRPCServer()   
    startHTTPServer()
}
