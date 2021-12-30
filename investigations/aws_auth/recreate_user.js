if (!process.env.hasOwnProperty("TEST_AWS_ARN")) {
    print ("TEST_AWS_ARN environment variable not set");
    print ("Run . ./get_temporary_credentials.sh first");
}

const db = db.getSiblingDB("$external");
try {
    db.dropUser(process.env.TEST_AWS_ARN);
} catch (e) {
    print ("Ignoring error on dropUser: " + e);
}

db.createUser ({
    user: process.env.TEST_AWS_ARN,
    roles: [ {role: "dbOwner", db: "db"} ]
})