import requests
import re
import tabulate

r = requests.get("https://api.github.com/repos/mongodb/mongo-go-driver/releases", headers={
    "Accept": "application/vnd.github.v3+json"
}, params={"per_page": "100"})

# print (r.text)

data = []
for release in r.json():
    # Skip patch releases.
    if not re.match(r"v[1-9]+\.[0-9]+\.0$", release["tag_name"]):
        continue
    data.append({
        "tag": "[{}]({})".format(release["tag_name"], release["url"]),
        "date": release["created_at"].split("T")[0]
    })

print (tabulate.tabulate (data, headers="keys", tablefmt="markdown"))