FROM scratch
EXPOSE 8080
ENTRYPOINT ["/metadata-mgmt-services"]
COPY ./bin/ /
COPY ./conf/ /conf