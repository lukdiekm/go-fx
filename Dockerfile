FROM docker
COPY go-fx /go-fx
RUN chmod +x /go-fx
CMD [ "/go-fx" ]