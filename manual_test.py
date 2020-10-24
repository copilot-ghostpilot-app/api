#!/usr/bin/env python3

import json
from requests import request

class API(object):

    @staticmethod
    def get_api(endpoint="http://0.0.0.0:8080/tweets/emojis", payload="Test"):
        headers= {}
        response = request("GET", endpoint, headers=headers, data=payload)
        return response

    @staticmethod
    def post_api(endpoint="http://0.0.0.0:8080/tweets/create", payload="Test"):
        headers = {'Content-Type': 'application/json'}
        response = request("POST", endpoint, headers=headers, data=payload)
        return response

if __name__ == "__main__":
    _payload = {'id': '1319721334776029184', 'username': 'efekarakus', 'tweet_content': ' Copilot v0.5.0 is released!! \U0001f973\n Manage scheduled jobs with the new "copilot job" cmds.\n Use an existing image with the new "location" field.\n Minor improvements: backend services no longer require ports, and deletions don\'t need a profile!\n\nhttps://t.co/fME5yysOPT', 'metadata': {'media': None, 'hashtags': [], 'created_date': 'Fri Oct 23 19:24:20 +0000 2020', 'retweet_data': None}}
    write_api = API.post_api(payload=json.dumps(_payload).encode("utf-8"))
    print(write_api)

    read_api = API.get_api()
    print(read_api)
