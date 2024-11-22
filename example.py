import requests


good_request= {"currentTimezone": "EST",
    "destinationTimezone": "PST",
}

r= requests.post('http://localhost:8080/time', json=good_request)
print(r.text)
