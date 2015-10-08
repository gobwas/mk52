# Electronica MK-52
___________________

![mk52](https://upload.wikimedia.org/wikipedia/commons/8/82/Elektronika_MK_52_with_accessories.jpg)

# Build

This app uses [gb](http://getgb.io) to manage dependencies.
So use it to restore deps and build sources:

```shell
    gb vendor restore && gb build
```

# Use

Both `client` and `server` apps use the same options:

Option | Default   | Meaning
-------|-----------|--------
host     localhost   Host to listen/connect
port     5555        Port to listen/connect
route    mk52        Route for websocket handler
timeout  10          Timeout in seconds to wait for request/response

So:

```shell
./bin/server
```

```shell
./bin/client - "1 + 1 + 5 * 10 + 9 ^ 2"
```

# Output

Both `client` and `server` use the same logging strategy - all log message is sent to `stderr`.
`client` send result data to `stdout`.