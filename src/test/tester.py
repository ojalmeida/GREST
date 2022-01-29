import json
import time

import requests
import sys

path_mappings_to_test = (
    {"path_mapping_id": "100", "path": "/authors", "table": "authors"}, {"path": "/posts", "table": "posts"})
key_mappings_to_test = (
    {"key_mapping_id": "200", "key": "name", "column": "first_name"}, {"key": "description", "column": "bio"})
behaviors_to_test = ({"behavior_id": "300", "path_mapping_id": "100", "key_mapping_id": "200"},
                     {"behavior_id": None, "path_mapping_id": None, "key_mapping_id": None})


# Missing keys, None values  --> Auto assigned by database


def test():
    sys.stdout.write("# Checking GREST server - ")
    sys.stdout.write("Ok\n\n") if checkUp() else sys.stderr.write("Error\n\n")

    sys.stdout.write("\n# Checking endpoints\n\n")

    time.sleep(1)
    test = test_path_mapping_endpoint()

    if test[0]:
        sys.stdout.write("\tpath-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tpath-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    time.sleep(1)
    test = test_key_mapping_endpoint()

    if test[0]:
        sys.stdout.write("\tkey-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tkey-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    time.sleep(1)
    test = test_behavior_endpoint()

    if test[0]:
        sys.stdout.write("\tbehavior -> Ok\n")

    else:
        sys.stderr.write(f"\n\tbehavior -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    sys.stdout.write("\n\n# Checking insertions\n\n")

    time.sleep(1)
    test = test_path_mapping_insertion()

    if test[0]:
        sys.stdout.write("\tpath-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tpath-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    time.sleep(1)
    test = test_key_mapping_insertion()

    if test[0]:
        sys.stdout.write("\tkey-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tkey-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    time.sleep(1)
    test = test_behavior_insertion()

    if test[0]:
        sys.stdout.write("\tbehavior -> Ok\n")

    else:
        sys.stderr.write(f"\n\tbehavior -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    sys.stdout.write("\n\n# Checking retrieving\n\n")

    time.sleep(1)
    test = test_path_mapping_retrieve()

    if test[0]:
        sys.stdout.write("\tpath-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tpath-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n\n")

    time.sleep(1)
    test = test_key_mapping_retrieve()

    if test[0]:
        sys.stdout.write("\tkey-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tkey-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n\n")

    time.sleep(1)
    test = test_behavior_retrieve()

    if test[0]:
        sys.stdout.write("\tbehavior -> Ok\n")

    else:
        sys.stderr.write(f"\n\tbehavior -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    sys.stdout.write("\n\n# Checking change\n")

    time.sleep(1)
    test = test_path_mapping_change()

    if test[0]:
        sys.stdout.write("\tpath-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tpath-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    time.sleep(1)
    test = test_key_mapping_change()

    if test[0]:
        sys.stdout.write("\tkey-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tkey-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    time.sleep(1)
    test = test_behavior_change()

    if test[0]:
        sys.stdout.write("\tbehavior -> Ok\n")

    else:
        sys.stderr.write(f"\n\tbehavior -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    sys.stdout.write("\n# Checking deletions\n\n")

    time.sleep(1)
    test = test_path_mapping_deletion()

    if test[0]:
        sys.stdout.write("\tpath-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tpath-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    time.sleep(1)
    test = test_key_mapping_deletion()

    if test[0]:
        sys.stdout.write("\tkey-mapping -> Ok\n")

    else:
        sys.stderr.write(f"\n\tkey-mapping -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")

    time.sleep(1)
    test = test_behavior_deletion()

    if test[0]:
        sys.stdout.write("\tbehavior -> Ok\n")

    else:
        sys.stderr.write(f"\n\tbehavior -> Error\n")

        sys.stderr.write("\t\t- request\n")
        sys.stderr.write(f"\t\t\t├──{test[1]['url']}\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['body']}\n")
        sys.stderr.write("\t\t- response\n")
        sys.stderr.write(f"\t\t\t└──{test[1]['response']}\n")


def checkUp() -> list:
    raw_response = requests.get(f"http://{grest_ip}:8080").text

    success = 0
    err = {}

    if raw_response == '404 page not found\n':
        success += 1

    return [success == 1, err]


def test_behavior_endpoint() -> list:
    success = 0
    err = {"url": None, "body": None, "response": None}

    url = f"http://{grest_ip}:9090/config/behaviors"

    raw_response = requests.get(url)

    response_obj = json.loads(raw_response.content)

    try:

        if response_obj["status"] == 200:
            success += 1

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_path_mapping_endpoint() -> list:
    success = 0
    err = {}

    raw_response = requests.get(f"http://{grest_ip}:9090/config/path-mappings")

    response_obj = json.loads(raw_response.content)

    try:

        if response_obj["status"] == 200:
            success += 1

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_key_mapping_endpoint() -> list:
    success = 0
    err = {}

    raw_response = requests.get(f"http://{grest_ip}:9090/config/behaviors")

    response_obj = json.loads(raw_response.content)

    try:

        if response_obj["status"] == 200:
            success += 1

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_behavior_insertion() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/behaviors"

    payload = {"Set": behaviors_to_test[0]}

    raw_response = requests.post(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception


    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    try:

        path_mapping_id = \
            json.loads(requests.get(f"http://{grest_ip}:9090/config/path-mappings", params={"table": "posts"}).content)[
                'response'][0]['path_mapping_id']
        key_mapping_id = \
            json.loads(requests.get(f"http://{grest_ip}:9090/config/key-mappings", params={"column": "bio"}).content)[
                'response'][0]['key_mapping_id']

        payload = {"Set": {"path_mapping_id": str(path_mapping_id), "key_mapping_id": str(key_mapping_id)}}

        raw_response = requests.post(url, data=json.dumps(payload))

        res = json.loads(raw_response.content)

        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_path_mapping_insertion() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/path-mappings"

    payload = {"Set": path_mappings_to_test[0]}

    raw_response = requests.post(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    payload = {"Set": path_mappings_to_test[1]}

    raw_response = requests.post(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_key_mapping_insertion() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/key-mappings"

    payload = {"Set": key_mappings_to_test[0]}

    raw_response = requests.post(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    payload = {"Set": key_mappings_to_test[1]}

    raw_response = requests.post(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_path_mapping_deletion() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/path-mappings"

    payload = {"Must": path_mappings_to_test[0]}

    raw_response = requests.delete(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    payload = {"Must": path_mappings_to_test[1]}

    raw_response = requests.delete(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_key_mapping_deletion() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/key-mappings"

    payload = {"Must": key_mappings_to_test[0]}

    raw_response = requests.delete(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    payload = {"Must": key_mappings_to_test[1]}

    raw_response = requests.delete(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_behavior_deletion() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/behaviors"

    payload = {"Must": behaviors_to_test[0]}

    raw_response = requests.delete(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.5
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    try:

        behavior_id = None

        for obj in json.loads(requests.get(url).content)['response']:

            if obj['behavior_id'] != 300:
                behavior_id = obj['behavior_id']
                break

        if behavior_id is not None:

            payload = {"Must": {"behavior_id": behavior_id}}

            raw_response = requests.delete(url, data=json.dumps(payload))

            res = json.loads(raw_response.content)

            if res['status'] == 200:
                success += 1
            else:
                raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_path_mapping_change() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/path-mappings"

    payload = {"Must": path_mappings_to_test[0], "Set": {"path": "author"}}

    raw_response = requests.put(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 1

        else:
            raise Exception


    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_key_mapping_change() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/key-mappings"

    payload = {"Must": key_mappings_to_test[0], "Set": {"key": "description"}}

    raw_response = requests.put(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 1

        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_behavior_change() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/behaviors"

    payload = {"Must": behaviors_to_test[0], "Set": {"behavior_id": "999"}}

    raw_response = requests.put(url, data=json.dumps(payload))

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 1

        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_behavior_retrieve() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/behaviors"

    params = {"key_mapping_id": behaviors_to_test[0]['key_mapping_id'],
              "path_mapping_id": behaviors_to_test[0]['path_mapping_id']}

    raw_response = requests.get(url, params=params)

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.25
            
        else:
            raise Exception


        if len(res['response']) == 1:
            success += 0.25
            
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    params = {"behavior_id": behaviors_to_test[0]['behavior_id']}

    raw_response = requests.get(url, params=params)

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.25
            
        else:
            raise Exception

        if len(res['response']) == 1:
            success += 0.25
            
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_path_mapping_retrieve() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/path-mappings"

    params = {"path_mapping_id": path_mappings_to_test[0]['path_mapping_id'], "path": path_mappings_to_test[0]['path'],
              "table": path_mappings_to_test[0]['table']}

    raw_response = requests.get(url, params=params)

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.25
            
        else:
            raise Exception

        if len(res['response']) == 1:
            success += 0.25
            
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    params = {"path": path_mappings_to_test[1]['path'], "table": path_mappings_to_test[1]['table']}

    raw_response = requests.get(url, params=params)

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.25
            
        else:
            raise Exception

        if len(res['response']) == 1:
            success += 0.25
            
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


def test_key_mapping_retrieve() -> list:
    success = 0
    err = {}

    url = f"http://{grest_ip}:9090/config/key-mappings"

    params = {"key_mapping_id": key_mappings_to_test[0]['key_mapping_id'], "key": key_mappings_to_test[0]['key'],
              "column": key_mappings_to_test[0]['column']}

    raw_response = requests.get(url, params=params)

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.25
            
        else:
            raise Exception

        if len(res['response']) == 1:
            success += 0.25
            
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    params = {"key": key_mappings_to_test[1]['key'], "column": key_mappings_to_test[1]['column']}

    raw_response = requests.get(url, params=params)

    res = json.loads(raw_response.content)

    try:
        if res['status'] == 200:
            success += 0.25
            
        else:
            raise Exception

        if len(res['response']) == 1:
            success += 0.25
            
        else:
            raise Exception

    except Exception:

        err["url"] = raw_response.request.url
        err["body"] = raw_response.request.body
        err["response"] = raw_response.content
        pass

    return [success == 1, err]


mysql_ip = ""
grest_ip = ""

if __name__ == '__main__':

    # Getting IPs
    next_is_grest = False
    next_is_mysql = False

    for arg in sys.argv:

        if next_is_grest:
            grest_ip = arg
            next_is_grest = False

        if next_is_mysql:
            mysql_ip = arg
            next_is_mysql = False

        if arg == "--grest":
            next_is_grest = True

        if arg == "--mysql":
            next_is_mysql = True

    test()
