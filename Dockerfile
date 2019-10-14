FROM scratch
EXPOSE 8080
ENTRYPOINT ["/meta-data-mgmt"]
COPY ./bin/ /