INSTANCE_NAME=binz
docker run -d \
    -e WAVEFRONT_URL=https://${INSTANCE_NAME}.wavefront.com/api/ \
    -e WAVEFRONT_TOKEN=${TO_TOKEN} \
    -e JAVA_HEAP_USAGE=512m \
    -e WAVEFRONT_PROXY_ARGS="--customTracingListenerPorts 30001" \
    -p 2878:2878 \
    -p 30001:30001 \
    wavefronthq/proxy:latest

