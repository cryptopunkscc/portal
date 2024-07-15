```mermaid
graph LR
    cli --> app
    cli --> dev
    app --> app-goja
    app --> app-wails
    app --> app-exec
    dev --> dev-goja
    dev --> dev-wails
    dev --> dev-go
    dev-goja --> goja
    dev-goja --> goja-dev
    dev-goja --> goja-dist
    dev-wails --> wails
    dev-wails --> wails-dev
    dev-wails --> wails-dist
    dev-go --> exec
```