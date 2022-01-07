The `certs` directory contains files constructed from the `x509gen` directory. The `x509gen` directory is the download from https://x509gen.corp.mongodb.com/#/cert/5ce5b21a42a0ef0008b11399. See `make_certs.py` for how the files are concatenated.

Use `openssl x509 -text -in <cert.pem> -noout` to inspect a .pem file.