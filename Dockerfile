FROM scratch
COPY gocd-metrics /gocd-metrics
EXPOSE 9090
ENTRYPOINT ["/gocd-metrics"]
