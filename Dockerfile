# Create a minimal container to run a Golang static binary
FROM scratch
COPY tictac /
ENTRYPOINT ["/tictac"]
EXPOSE 8080
