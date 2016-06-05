FROM debian

WORKDIR /root

ADD p3ka ./p3ka
ADD static static

ENTRYPOINT ["./p3ka"]