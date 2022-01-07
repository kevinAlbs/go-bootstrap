# Combine certs from https://x509gen.corp.mongodb.com/#/cert/5ce5b21a42a0ef0008b11399
certs = ["./x509gen/Drivers-Testing-Client-Second-Level.key",
         "./x509gen/Drivers-Testing-Client-Second-Level.pem",
         "./x509gen/Drivers-Testing-Client-Intermediate.pem"]

with open("./certs/Drivers-Testing-Client-Intermediate-combined.pem", "w") as outfile:
    for cert in certs:
        with open(cert, "r") as infile:
            outfile.write(infile.read())
            outfile.write("\n")

certs = ["./x509gen/Drivers-Testing-Server.key",
         "./x509gen/Drivers-Testing-Server.pem"]

with open("./certs/Drivers-Testing-Server-combined.pem", "w") as outfile:
    for cert in certs:
        with open(cert, "r") as infile:
            outfile.write(infile.read())
            outfile.write("\n")