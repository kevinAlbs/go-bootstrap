from pymongo import MongoClient
from bson.binary import UuidRepresentation
from uuid import uuid4

# use the 'standard' representation for cross-language compatibility.
client = MongoClient(uuidRepresentation="pythonLegacy")
collection = client.get_database('uuid_db').get_collection('uuid_coll')

# remove all documents from collection
collection.delete_many({})

# create a native uuid object
uuid_obj = uuid4()

# save the native uuid object to MongoDB
collection.insert_one({'uuid': uuid_obj})