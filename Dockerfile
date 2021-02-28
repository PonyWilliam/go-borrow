FROM alpine
ADD borrow /borrow
ENTRYPOINT [ "/borrow" ]
