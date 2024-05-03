```mermaid
flowchart
    cli.Open --> open.Run
    open.Run -- is a service --> portal.Cmd
    open.Run -- not a service --> portal.Open
    portal.Open --> portal.open
    portal.open -. OK .-> open.Run
    portal.open -- EOF --> serve.Run
    serve.Run --> open.Run
```