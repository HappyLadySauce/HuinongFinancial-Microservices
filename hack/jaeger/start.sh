#!/bin/bash
docker run --rm --name jaeger -p 16686:16686 -p 4317:4317 -p 4318:4318 -p 5778:5778  -p 9411:9411 -p 14268:14268 jaegertracing/all-in-one:latest