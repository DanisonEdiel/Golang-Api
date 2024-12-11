package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/rpc"
)

type Args struct {
    A, B int
}

type Response struct {
    Result int `json:"result"`
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Conecta al servidor RPC
    client, err := rpc.Dial("tcp", "localhost:1234")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error connecting to RPC server: %v", err), http.StatusInternalServerError)
        return
    }
    defer client.Close()

    // Lee los parámetros de la solicitud
    var args Args
    if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
        http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
        return
    }

    fmt.Printf("Received RPC request: %+v\n", args) // Log para depurar

    // Llama al método Multiply del servidor RPC
    var result int
    err = client.Call("Calculator.Multiply", &args, &result)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error calling RPC method: %v", err), http.StatusInternalServerError)
        return
    }

    // Responde al frontend con el resultado
    response := Response{Result: result}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func main() {
    // Define el endpoint HTTP
    http.HandleFunc("/rpc/multiply", rpcHandler)

    fmt.Println("HTTP server listening on port 8081")
    http.ListenAndServe(":8081", nil)
}
