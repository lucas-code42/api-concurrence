import requests
from concurrent.futures import ThreadPoolExecutor

url = "http://localhost:8080/01000001"
total_requests = 2000
concurrent_requests = 1000

def make_request(_):
    response = requests.get(url)
    print(response.json())

    if response.status_code != 200:
        print(f"status_code: {response.status_code}")

with ThreadPoolExecutor(max_workers=concurrent_requests) as executor:
    executor.map(make_request, range(total_requests))

print("Done!")
