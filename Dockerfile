FROM gcr.io/distroless/static:nonroot

# `nonroot` coming from distroless
USER 65532:65532

# Copy the binary that goreleaser built
COPY  load-balancer-api /load-balancer-api

# Run the web service on container startup.
ENTRYPOINT ["/load-balancer-api"]
CMD ["serve"]
