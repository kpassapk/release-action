# File generated by github.com/posener/goaction. DO NOT EDIT.


FROM 1.20.3-alpine3.17
RUN apk add git 

COPY . /home/src
WORKDIR /home/src
RUN go build -o /bin/action ./

ENTRYPOINT [ "/bin/action" ]
