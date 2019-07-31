FROM gcr.io/distroless/static
COPY gocd-metrics /gocd-metrics
EXPOSE 9090
ENTRYPOINT ["/gocd-metrics"]
