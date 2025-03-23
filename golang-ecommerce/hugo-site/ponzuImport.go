package main

import (
    "bufio"
    "crypto/tls"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "time"
)

type PonzuProductResponse struct {
    Data []PonzuProduct `json:"data"`
}

type PonzuProduct struct {
    UUID        string  `json:"uuid"`
    ID          int     `json:"id"`
    Slug        string  `json:"slug"`
    Timestamp   int64   `json:"timestamp"`
    Updated     int64   `json:"updated"`
    Name        string  `json:"name"`
    Price       float32 `json:"price"`
    Description string  `json:"description"`
    Image       string  `json:"image"`
}

type SnipcartProductResponse struct {
    Items []SnipcartProduct `json:"items"`
}

type SnipcartProduct struct {
    Stock int `json:"stock"`
}

type HugoProduct struct {
    ID               string    `json:"id"`
    Title            string    `json:"title"`
    Date             time.Time `json:"date"`
    LastModification time.Time `json:"lastmod"`
    Description      string    `json:"description"`
    Price            float32   `json:"price"`
    Image            string    `json:"image"`
    Stock            int       `json:"stock"`
}

func (dest *HugoProduct) mapPonzuProduct(src PonzuProduct, ponzuHostURL string, client *http.Client) {
    dest.ID = src.Slug
    dest.Title = src.Name
    dest.Price = src.Price
    dest.Description = src.Description
    dest.Image = ponzuHostURL + src.Image
    dest.Date = time.Unix(src.Timestamp/1000, 0)
    dest.LastModification = time.Unix(src.Updated/1000, 0)

    url := "https://app.snipcart.com/api/products?userDefinedId=" + dest.ID
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        log.Printf("Failed to create Snipcart request: %v", err)
        return
    }
    apiKey := base64.StdEncoding.EncodeToString([]byte(os.Getenv("SNIPCART_PRIVATE_API_KEY")))
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", "Basic "+apiKey)
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Failed to fetch Snipcart stock: %v", err)
        return
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Failed to read Snipcart response: %v", err)
        return
    }
    var products SnipcartProductResponse
    if err = json.Unmarshal(body, &products); err != nil {
        log.Printf("Failed to unmarshal Snipcart response: %v", err)
        return
    }
    if len(products.Items) > 0 {
        dest.Stock = products.Items[0].Stock
    }
}

func main() {
    ponzuHostURL := os.Getenv("PONZU_HOST_URL")
    if ponzuHostURL == "" {
        log.Fatal("PONZU_HOST_URL environment variable not set")
    }
    productsEndpoint := ponzuHostURL + "/api/contents?type=Product"

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Remove in production
    }
    client := &http.Client{Transport: tr}

    resp, err := client.Get(productsEndpoint)
    if err != nil {
        log.Fatalf("Failed to fetch Ponzu products: %v", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read Ponzu response: %v", err)
    }
    var products PonzuProductResponse
    if err = json.Unmarshal(body, &products); err != nil {
        log.Fatalf("Failed to unmarshal Ponzu response: %v", err)
    }

    dir, err := os.Getwd()
    if err != nil {
        log.Fatalf("Failed to get working directory: %v", err)
    }
    contentDir := dir + "/content/product"
    if err := os.RemoveAll(contentDir); err != nil {
        log.Fatalf("Failed to remove content directory: %v", err)
    }
    if err := os.MkdirAll(contentDir, 0755); err != nil {
        log.Fatalf("Failed to create content directory: %v", err)
    }

    for _, srcProduct := range products.Data {
        var destProduct HugoProduct
        destProduct.mapPonzuProduct(srcProduct, ponzuHostURL, client)
        filePath := fmt.Sprintf("%s/%s.md", contentDir, destProduct.ID)
        file, err := os.Create(filePath)
        if err != nil {
            log.Printf("Failed to create file %s: %v", filePath, err)
            continue
        }
        writer := bufio.NewWriter(file)
        frontMatter, _ := json.MarshalIndent(destProduct, "", "  ")
        writer.WriteString("---\n")
        writer.WriteString(string(frontMatter) + "\n")
        writer.WriteString("---\n")
        writer.WriteString(destProduct.Description + "\n")
        writer.Flush()
        file.Close()
    }
}