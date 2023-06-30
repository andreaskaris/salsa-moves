FROM registry.fedoraproject.org/fedora-minimal:latest AS build
RUN microdnf install golang libX*-devel libglvnd-devel -y && microdnf clean all -y
WORKDIR /app
COPY . .
RUN make build

FROM registry.fedoraproject.org/fedora-minimal:latest AS build-release
RUN microdnf install golang libX*-devel libglvnd-devel -y && microdnf clean all -y
WORKDIR /
COPY --from=build /app/_output/salsa-moves /usr/local/bin/salsa-moves
COPY --from=build /app/config.yaml /config.yaml
