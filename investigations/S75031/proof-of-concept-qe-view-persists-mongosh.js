// Refer: https://gist.github.com/ZacharyEspiritu/31dc98ed6bc525ca75f490600d0109f5
// Run with: ~/bin/mongosh-1.8.0-darwin-arm64/bin/mongosh ./proof-of-concept-qe-view-persists-mongosh.js
/**
 * This file is a proof of concept for mongosh showing that an
 * encrypted client can create a view over a QE-encrypted field
 * (with no error) and an unencrypted client can see the view
 * definition on the server in plaintext. It does the following:
 *
 *  1. Creates a Mongo client with automatic encryption options.
 *  2. Uses the encrypted client to create a QE-encrypted collection
 *     with a QE-encrypted field.
 *  3. Creates a view using the encrypted client on the QE-encrypted
 *     collection with a $match expression over the QE-encrypted
 *     field.
 *  4. Uses an *unencrypted* Mongo client to view the system.views
 *     collection.
 *
 * Tested using combination of:
 *
 *  - mongosh v1.8.0
 *  - 10gen/mongo v7.0.0-alpha commit 82dc951
 */

// URI of the mongo instance to connect to.
const MONGO_URI = "mongodb://127.0.0.1:27017";

// Constants for the proof of concept.
const QE_DATABASE_NAME = "proofOfConceptForMongoShell";
const QE_COLLECTION_NAME = "myQECollection";

// Set up an unencrypted client and database.
const unencryptedClient = Mongo(MONGO_URI);
const unencryptedDb = unencryptedClient.getDB(QE_DATABASE_NAME);

// Drop the database in case you've run the proof of concept before.
unencryptedDb.dropDatabase();

// Set up an encrypted client and database with a test key.
const autoEncryptionOpts = {
    keyVaultNamespace: `${QE_DATABASE_NAME}.keystore`,
    kmsProviders: {
        local: {
            key: BinData(0, "/tu9jUCBqZdwCelwE/EAm/4WqdxrSMi04B8e9uAV+m30rI1J2nhKZZtQjdvsSCwuI4erR6IEcEK+5eGUAODv43NDNIR9QheT2edWFewUfHKsl9cnzTc86meIzOmYl6dr")
        }
    }
};
const encryptedClient = Mongo(MONGO_URI, autoEncryptionOpts);
const encryptedDb = encryptedClient.getDB(QE_DATABASE_NAME);

// Create a collection with an encrypted "myEncryptedField" field 
// with some test documents.
const keyVault = encryptedClient.getKeyVault();
keyVault.createKey(
    "local", ['encryptedFieldKey']);
encryptedDb.createCollection(QE_COLLECTION_NAME, {
    encryptedFields: {
        fields: [{
            path: "myEncryptedField",
            keyId: keyVault.getKeyByAltName("encryptedFieldKey")._id,
            bsonType: "int",
            queries: {
                "queryType": "equality"
            },
        },],
    },
});
encryptedDb.runCommand({
    "insert": QE_COLLECTION_NAME,
    "documents": [{
        myEncryptedField: NumberInt(1)
    }, {
        myEncryptedField: NumberInt(2)
    },],
});

// Use the encrypted client to create a view over the collection
// with the QE-encrypted field.
encryptedDb.createCollection("myView", {
    viewOn: QE_COLLECTION_NAME,
    pipeline: [{
        $match: {
            myEncryptedField: NumberInt(3)
        }
    }]
});

// Attempt to see the view using the unencrypted client.
const result = unencryptedDb.system.views.find();
print(`Result of unencrypted client querying ${QE_DATABASE_NAME}.system.views:`);
printjson(result);

// Outputs:
// [
//   {
//     _id: `${QE_DATABASE_NAME}.myView`,
//     viewOn: QE_COLLECTION_NAME,
//     pipeline: [ { '$match': { myEncryptedField: 3 } } ]
//   }
// ]