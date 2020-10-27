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
    _payload = {"id": "1320825299236325122", "username": "bvtujo", "tweet_content": "So fun to hang out with Adam and @brentContained \U0001f9c3 talking about the new Scheduled Jobs feature in #awscopilot \ud83d\udc7b\u2620\ufe0f\ud83e\udd16\ud83e\udd20\ud83c\udf83 https://t.co/MPmXuDg3Xn", "metadata": {"media": "None", "hashtags": "[Hashtag(Text='awscopilot')]", "created_date": "Mon Oct 26 20:31:05 +0000 2020", "retweet_data": "None"}}
    write_api = API.post_api(payload=json.dumps(_payload).encode("utf-8"))
    print(write_api)

    read_api = API.get_api()
    print(read_api)
