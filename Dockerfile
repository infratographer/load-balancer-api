FROM gcr.io/distroless/static:nonroot

# `nonroot` coming from distroless
USER 65532:65532

# Copy the binary that goreleaser built
COPY  loadbalancer-api /loadbalancer-api

# Run the web service on container startup.
ENTRYPOINT ["/loadbalancer-api"]
CMD ["serve"]
