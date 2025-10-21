FROM scratch

COPY build/dist-sensor-if /dist-sensor-if
COPY static /static

WORKDIR /

EXPOSE 80

ENTRYPOINT ["/dist-sensor-if", "-addr", ":80"]
