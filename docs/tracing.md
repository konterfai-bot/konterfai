[<- back to docs](README.md)

# Tracing

## Setting up tracing

To set up tracing, you need to have a Jaeger instance, or any other opentelemetry client running. 
You can run jaeger using Docker:

```bash
$> docker run --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  --restart always \
  jaegertracing/all-in-one:1.35
```

The web UI will be available at [http://localhost:16686](http://localhost:16686).

To enable tracing in konterfAI, you need to use the cli-flag `--tracing-endpoint`:

```bash
$> konterfai --tracing-endpoint localhost:4317
```

This option can be combined with other flags.