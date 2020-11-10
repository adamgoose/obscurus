FROM vault

COPY ./obscurus.hcl /obscurus.hcl
COPY ./vault-setup.sh /vault-setup.sh

ENTRYPOINT [ "/vault-setup.sh" ]