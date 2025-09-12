Used to manually verify Rust driver can spawn mongocryptd on Windows:

- Spawned an Evergreen a `windows-2022-latest-large` host
- Installed libmongocrypt as described [here](https://pkg.go.dev/go.mongodb.org/mongo-driver/v2@v2.3.0/mongo#hdr-In_Use_Encryption).
- Downloaded [latest Windows build](https://downloads.10gen.com/windows/mongodb-windows-x86_64-enterprise-latest.zip) (8.3.0-alpha0-1179-gb3a8643)
- Ensured no mongocryptd processes are running in Powershell: `Get-Process -Name "mongocryptd"`
- Ran `test.sh` from Cygwin to test auto encryption triggering mongocryptd spawning:
```
$ ./test.sh
Test passed!
```
- Ensured mongocryptd was spawned by running: `Get-Process -Name "mongocryptd"`


