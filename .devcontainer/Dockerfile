FROM ubuntu:latest

# Install necessary packages
RUN apt-get update && \
    apt-get install -y curl git vim zsh sudo locales tree software-properties-common && \
    rm -rf /var/lib/apt/lists/* # Clean up apt cache to reduce image size

# Set timezone to JST
ENV TZ=Asia/Tokyo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Set locale to Japanese
RUN locale-gen ja_JP.UTF-8
ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8

# Install Golang from PPA
# Add -y to avoid interactive prompts
RUN add-apt-repository ppa:longsleep/golang-backports -y && \
    apt-get update && \
    apt-get install -y golang-go && \
    rm -rf /var/lib/apt/lists/* # Clean up apt cache

# Set Go environment variables
# Use /go as GOPATH, which is the standard default in recent Go versions
ENV GOPATH /go
# Ensure standard Go installation path and GOPATH/bin are in PATH
ENV PATH $GOPATH/bin:/usr/local/go/bin:$GOROOT/bin:$PATH
ENV GO111MODULE on

# Install Cobra CLI tool globally
# This tool helps generate Cobra application structure
RUN go install github.com/spf13/cobra-cli@latest

# Set working directory
WORKDIR /workspace

# Set default command
# Note: docker-compose service command in compose.yml overrides this CMD
CMD [ "/bin/zsh" ]