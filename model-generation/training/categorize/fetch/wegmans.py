"""Module for fetching Wegman's data

Notes:
1. Make a regular HTTP request to wegmans.com to get a valid cookie session. This is looked for by the backend server
2. Use their REST API (https://shop.wegmans.com/api/v2/categories?store_id=7) to request departments
  a. note that it requires a store ID. For training purpsoes, we want to largest selection
3. Use thier REST API (https://shop.wegmans.com/api/v2/store_products?category_id=144&category_ids=144&limit=60&offset=0&sort=popular
Request Method: GET)) to request items in each department
4. Map to standard departments
"""
import requests
import logging

# These two lines enable debugging at httplib level (requests->urllib3->http.client)
# You will see the REQUEST, including HEADERS and DATA, and RESPONSE with HEADERS but without DATA.
# The only thing missing will be the response.body which is not logged.
try:
    import http.client as http_client
except ImportError:
    # Python 2
    import httplib as http_client
http_client.HTTPConnection.debuglevel = 1
# You must initialize logging, otherwise you'll not see debug output.
logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)

def get_wegmans_session() -> requests.Session:

    session = requests.Session()
    r = session.post(
        "https://shop.wegmans.com/api/v2/user_sessions",
        json={
            "binary": "web-ecom",
            "binary_version": "4.1.33",
            "is_retina": False,
            "os_version": "MacIntel",
            "pixel_density": "2.0",
            "push_token": "",
            "screen_height": 900,
            "screen_width": 1440,
        },
    )
    r.raise_for_status()

    return session


def make_api_call(session: requests.Session):
    headers = {
      "Connection": "keep-alive",
      "Accept": "*/*",
      "Accept-Encoding": "gzip, deflate, br",
      "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
    }
    params = {
      "category_id": 144,
      "category_ids": 144,
      "limit": 60,
      "offset": 0,
      "sort": "popular"
    }

    URL = "https://shop.wegmans.com/api/v2/store_products"
    print("")
    print(f"Cookies: {session.cookies}")
    print("")
    resp = requests.get(URL, params=params, headers=headers, cookies=session.cookies)

    resp.raise_for_status()
    return resp.json()


if __name__ == "__main__":
  session = get_wegmans_session()
  j = make_api_call(session)
  print(j)
