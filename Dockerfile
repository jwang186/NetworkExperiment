FROM envoyproxy/envoy:v1.21-latest
COPY demo.yaml /etc/envoy/envoy.yaml
RUN chmod go+r /etc/envoy/envoy.yaml
