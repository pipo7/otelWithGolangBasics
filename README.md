# OTEL with Go example 

refrence: https://blog.devgenius.io/working-with-opentelemetry-and-golang-5037ca301bc5
Oficial doc: https://www.jaegertracing.io/docs/1.45/

Jaeger provides a distribution all-in-one
For production secure the container , volume and access.
```
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.45
```

After running and building the container you can access Jaeger in the following link ```http://127.0.0.1:16686/search.  
NOTE: // Replace this with system's IP ```

# packages needed
"go.opentelemetry.io/otel"
"go.opentelemetry.io/otel/exporters/jaeger"

Configure the application to send the telemetry data to Jaeger. This function creates the exporter using the default Jaeger URL:PORT.

# Example1 basic go program with OTEL
Run the program to see traces
``` go run main.go & ```
If we run the program and then go to the Jaeger main page we can search by our stats by selecting the following params:
service: medium-tutorial
operations: all

Ensure all ports are open.

# Example2 with gin webserver with OTEL
```go run gin-gonic-otel.go &```
```curl http://localhost:8080/ping```
Then see the traces in JaegerUI 