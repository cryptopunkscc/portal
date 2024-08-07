namespace Rpc {

    interface Client extends Binder<Client>, Caller {
        targetId: string
        target(id: string): Client
        serve(ctx: any): Promise<void>
    }

    interface Binder<T> {
        bind<R>(...methods: string[]): T & R
    }

    interface Caller {
        call(port: string, ...args: any): Call
        conn(port: string, ...args: any): Promise<Conn>
    }

    interface Call extends Consumer {
        conn(): Promise<Conn>
    }

    interface Conn extends Binder<Conn>, Caller, Consumer {
        encode(data: any): Promise<void>
        decode<T>(): Promise<T>
        close(): Promise<void>
    }

    interface Consumer extends Single {
        query: string
        port: string
        params: any[]
        map<T, R>(f: Map): this
        request<T, R>(...args: any): Promise<R>
        collect<T, R>(...args: any): Promise<R[] | undefined>
    }

    type Map = <K, V>(key: K) => V | null

    type Single = <R>() => R | null
}
