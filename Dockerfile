FROM alpine:3.16.2 AS build
ENV WOL_VERSION 0.2.4
RUN apk add go git

# Install Wol
RUN git clone -b v${WOL_VERSION} --depth 1 https://github.com/tikobus/wol
RUN cd wol && go build .

FROM alpine:3.16.2
COPY --from=build /wol/wol /usr/local/bin/wol
CMD ["wol"]
