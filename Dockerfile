FROM ubuntu:22.04
RUN apt update && apt install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /app
ADD email-alert-processor ./
EXPOSE 8085/tcp
ENTRYPOINT [ "/bin/bash", "-c", "./email-alert-processor" ]