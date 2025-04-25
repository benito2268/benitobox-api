import requests
from datetime import datetime

url = "http://localhost/post"
data = {
    "author" : "ben",
    "datetime" : 1234,
    "content" : "test"
}

response = requests.post(url=url, json=data, timeout=60)
print(response)