FROM scratch
MAINTAINER Rid <rid@cylo.io>
ADD dist/hapesay_linux_amd64/hapesay hapesay
ENTRYPOINT ["/hapesay"]