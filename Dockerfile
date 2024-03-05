FROM ubuntu:22.04
WORKDIR /app
ADD email-alert-processor ./
EXPOSE 8085/tcp
ENTRYPOINT [ "/bin/bash", "-c", "./email-alert-processor" ]