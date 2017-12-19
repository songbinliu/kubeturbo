FROM alpine:3.3

MAINTAINER Enlin Xu <enlin.xu@turbonomic.com>

RUN apk --update upgrade && apk add ca-certificates && update-ca-certificates
COPY ./_output/kubeturbo.linux /bin/kubeturbo 
RUN chmod +x /bin/kubeturbo

ENTRYPOINT ["/bin/kubeturbo"]



