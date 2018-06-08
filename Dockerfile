FROM docker:17.12.0-ce-dind

ADD docker-publish /bin/
ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/docker-publish"]
