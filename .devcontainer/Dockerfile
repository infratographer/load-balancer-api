# Used to install CRDB into the devcontainer
FROM cockroachdb/cockroach:latest-v22.2 as CRDB

FROM mcr.microsoft.com/vscode/devcontainers/go:0-1.20-bullseye

# Set up crdb
RUN mkdir /usr/local/lib/cockroach
COPY --from=CRDB /cockroach/cockroach /usr/local/bin
COPY --from=CRDB /usr/local/lib/cockroach/libgeos.so /usr/local/lib/cockroach/
COPY --from=CRDB /usr/local/lib/cockroach/libgeos_c.so /usr/local/lib/cockroach/

# [Optional] Uncomment this section to install additional OS packages.
RUN export DEBIAN_FRONTEND=noninteractive \
    && curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg \
    &&  curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | gpg --dearmor > /usr/share/keyrings/yarn-archive-keyring.gpg \
    && apt-get install apt-transport-https --yes \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list \
    && apt-get update \
    && apt-get -y install --no-install-recommends \
        bash-completion \
        helm \
        uuid-runtime \
        postgresql-client


# Install NATS Tooling
RUN curl -o /tmp/install.sh  https://raw.githubusercontent.com/nats-io/nsc/main/install.sh \
    && chmod +x /tmp/install.sh \
    && su vscode -c "/tmp/install.sh -s /usr/local/bin" \
    && rm -f /tmp/install.sh \
# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode

ENV NATS_CLI_VERSION=0.0.35

RUN go install github.com/volatiletech/sqlboiler/v4@latest \
    && go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest \
    && go install github.com/glerchundi/sqlboiler-crdb/v4@latest \
    && go install github.com/nats-io/natscli/nats@v${NATS_CLI_VERSION} \
    && go install github.com/nats-io/nkeys/nk@latest

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

WORKDIR /workspace
