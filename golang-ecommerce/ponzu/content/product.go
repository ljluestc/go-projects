package content

import (
    "log"
    "net/http"
    "os"

    "github.com/ponzu-cms/ponzu/management/editor"
    "github.com/ponzu-cms/ponzu/system/item"
)

type Product struct {
    item.Item

    Name        string  `json:"name"`
    Price       float32 `json:"price"`
    Description string  `json:"description"`
    Image       string  `json:"image"`
}

func (p *Product) MarshalEditor() ([]byte, error) {
    view, err := editor.Form(
        editor.Field{
            View: editor.Input("Name", p, map[string]string{
                "label":       "Name",
                "placeholder": "Enter product name",
            }),
        },
        editor.Field{
            View: editor.Input("Price", p, map[string]string{
                "label":       "Price",
                "placeholder": "Enter price in dollars",
                "type":        "number",
                "step":        "0.01",
            }),
        },
        editor.Field{
            View: editor.Textarea("Description", p, map[string]string{
                "label":       "Description",
                "placeholder": "Enter product description",
            }),
        },
        editor.Field{
            View: editor.File("Image", p, map[string]string{
                "label": "Image",
            }),
        },
    )

    if err != nil {
        return nil, err
    }
    return view, nil
}

func (p *Product) AfterAdminCreate(res http.ResponseWriter, req *http.Request) error {
    sendWebHook()
    return nil
}

func (p *Product) AfterAdminUpdate(res http.ResponseWriter, req *http.Request) error {
    sendWebHook()
    return nil
}

func (p *Product) AfterAdminDelete(res http.ResponseWriter, req *http.Request) error {
    sendWebHook()
    return nil
}

func sendWebHook() {
    url := os.Getenv("NETLIFY_BUILD_HOOK_URL")
    if url == "" {
        log.Println("NETLIFY_BUILD_HOOK_URL not set, skipping webhook")
        return
    }
    resp, err := http.Post(url, "application/json", nil)
    if err != nil {
        log.Printf("Failed to call webhook: %v", err)
    } else {
        log.Printf("Webhook called successfully at %s with result %s", url, resp.Status)
        resp.Body.Close()
    }
}