// Test behavior of min/max/precision on collection creation.
// Run with: `mongosh test_create_fle2_range.js`.

// assertThrows asserts `fn` throws an exception containing the message `expectMessage` as a substring.
function assertThrows(fn, expectMessage) {
    let thrown = false;
    let gotMessage;

    try {
        fn();
    } catch (e) {
        thrown = true;
        gotMessage = e.message;
    }
    if (!thrown) {
        throw "Expected exception to be thrown, but got no exception";
    }
    if (gotMessage.indexOf(expectMessage)) {
        throw "Expected exception to contain message '" + expectMessage + "' but got '" + gotMessage + "'";
    }
}

function reset() {
    db.coll.drop();
}

// Not specifying either min/max for an int results in an error.
{
    reset();
    assertThrows(() => {
        db.runCommand({
            "create": "coll",
            "encryptedFields": {
                "fields": [
                    {
                        "keyId": UUID(),
                        "path": "encryptedIndexed",
                        "bsonType": "int",
                        "queries": {
                            "queryType": "rangePreview",
                            "contention": 0,
                            "sparsity": 1,
                        },
                    },
                ]
            }
        });
    }, "The field 'min' is missing but required for range index");
}

// Specifying only min for an int results in an error.
{
    reset();
    assertThrows(() => {
        db.runCommand({
            "create": "coll",
            "encryptedFields": {
                "fields": [
                    {
                        "keyId": UUID(),
                        "path": "encryptedIndexed",
                        "bsonType": "int",
                        "queries": {
                            "queryType": "rangePreview",
                            "contention": 0,
                            "sparsity": 1,
                            "min": 0
                        },
                    },
                ]
            }
        });
    }, "The field 'max' is missing but required for range index");
}

// Specifying only max for an int results in an error.
{
    reset();
    assertThrows(() => {
        db.runCommand({
            "create": "coll",
            "encryptedFields": {
                "fields": [
                    {
                        "keyId": UUID(),
                        "path": "encryptedIndexed",
                        "bsonType": "int",
                        "queries": {
                            "queryType": "rangePreview",
                            "contention": 0,
                            "sparsity": 1,
                            "max": 0
                        },
                    },
                ]
            }
        });
    }, "The field 'min' is missing but required for range index");
}

// Specifying both min and max for an int results in an error.
{
    reset();
    assertThrows(() => {
        db.runCommand({
            "create": "coll",
            "encryptedFields": {
                "fields": [
                    {
                        "keyId": UUID(),
                        "path": "encryptedIndexed",
                        "bsonType": "int",
                        "queries": {
                            "queryType": "rangePreview",
                            "contention": 0,
                            "sparsity": 1,
                            "max": 0
                        },
                    },
                ]
            }
        });
    }, "The field 'min' is missing but required for range index");
}

// Specifying no min/max/precision for a double is OK.
{
    reset()
    let res = db.runCommand({
        "create": "coll",
        "encryptedFields": {
            "fields": [
                {
                    "keyId": UUID(),
                    "path": "encryptedIndexed",
                    "bsonType": "double",
                    "queries": {
                        "queryType": "rangePreview",
                        "contention": 0,
                        "sparsity": 1
                    },
                },
            ]
        }
    });
    assert(res.ok == 1);
}

// Specifying only min for a double is an error.
{
    reset()
    assertThrows(() => {
        db.runCommand({
            "create": "coll",
            "encryptedFields": {
                "fields": [
                    {
                        "keyId": UUID(),
                        "path": "encryptedIndexed",
                        "bsonType": "double",
                        "queries": {
                            "queryType": "rangePreview",
                            "contention": 0,
                            "sparsity": 1,
                            "min": 0
                        },
                    },
                ]
            }
        });
    }, "Precision, min, and max must all be specified together for floating point fields");
}

console.log("Tests passed")
