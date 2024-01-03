package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "strings"
    "time"
    "io/ioutil"
)

type PoolConfig struct {
    Domains          [][]interface{}            `json:"domains"`
    PathBasedDomains map[string][][]interface{} `json:"path_based_domains"`
    CustomHeaders    map[string]string          `json:"custom_headers"`
}

var (
    domainPools     map[string]PoolConfig
)

func loadConfig(filePath string) (map[string]PoolConfig, error) {
    var config map[string]PoolConfig

    configFile, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(configFile, &config)
    if err != nil {
        return nil, err
    }

    return config, nil
}

func init() {
    log.SetFlags(0)
    rand.Seed(time.Now().UnixNano())

    var err error
    domainPools, err = loadConfig("redirect-config.json")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
}

func performRedirection(w http.ResponseWriter, r *http.Request) {
    pathParts := strings.SplitN(r.URL.Path, "/", 4)
    if len(pathParts) < 4 {
        http.Error(w, "Invalid URL format", http.StatusBadRequest)
        return
    }

    poolID := pathParts[2]
    path := pathParts[3]
    clientIP := r.RemoteAddr

    logData := constructLogData(poolID, path, clientIP)

    poolConfig, exists := domainPools[poolID]
    if !exists {
        logError(logData, "Pool not found")
        http.Error(w, "Pool not found", http.StatusNotFound)
        return
    }

    selectedDomains := poolConfig.Domains
    for pathKey, domains := range poolConfig.PathBasedDomains {
        if strings.HasPrefix(path, pathKey) {
            selectedDomains = domains
            break
        }
    }

    if len(selectedDomains) == 0 {
        logError(logData, "No domains available for redirection")
        http.Error(w, "No domains available for redirection", http.StatusNotFound)
        return
    }

    domain := chooseDomain(selectedDomains)
    newURL := fmt.Sprintf("https://%s/%s", domain, path)
    if query := r.URL.RawQuery; query != "" {
        newURL += "?" + query
    }

    for header, value := range poolConfig.CustomHeaders {
        w.Header().Set(header, value)
    }

    logData["event"] = "redirect"
    logData["redirected_to"] = newURL
    logData["custom_headers"] = logCustomHeaders(poolConfig.CustomHeaders)
    logJSON(logData)

    http.Redirect(w, r, newURL, http.StatusFound)
}

func chooseDomain(domains [][]interface{}) string {
    var totalWeight int
    for _, domain := range domains {
        weight, ok := domain[1].(float64)
        if !ok {
            log.Println("Error: Invalid weight type")
            return ""
        }
        totalWeight += int(weight)
    }

    randWeight := rand.Intn(totalWeight)
    for _, domain := range domains {
        weight, _ := domain[1].(float64)
        intWeight := int(weight)
        if randWeight < intWeight {
            return domain[0].(string)
        }
        randWeight -= intWeight
    }
    return ""
}

func constructLogData(poolID, path, clientIP string) map[string]string {
    return map[string]string{
        "pool_id":        poolID,
        "requested_path": path,
        "client_ip":      clientIP,
        "datetime":       time.Now().UTC().Format(time.RFC3339),
    }
}

func logError(data map[string]string, message string) {
    data["event"] = "error"
    data["error_message"] = message
    logJSON(data)
}

func logJSON(data map[string]string) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Println("Error marshalling log data:", err)
        return
    }
    log.Println(string(jsonData))
}

func logCustomHeaders(customHeaders map[string]string) string {
    var headerStrings []string
    for header, value := range customHeaders {
        headerStrings = append(headerStrings, fmt.Sprintf("%s: %s", header, value))
    }
    return strings.Join(headerStrings, ", ")
}

func loaderIOHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("loaderio-4cf2b9e6e5bfe74c32ffdf9796cd8e5b"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("healthy"))
}

func main() {
    http.HandleFunc("/redirect/", performRedirection)
    http.HandleFunc("/loaderio-4cf2b9e6e5bfe74c32ffdf9796cd8e5b/", loaderIOHandler)
    http.HandleFunc("/health/", healthHandler)
    log.Fatal(http.ListenAndServe(":80", nil))
}
