import pymongo

client = pymongo.MongoClient()
client["db"].command({"ping": 1})
client.close()