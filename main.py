import requests
import sys
import os


def login(session, username, password):
    print("Try to login...")
    headers = {
        "Accept": "application/json",
        "Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
        "Connection": "keep-alive",
        "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
        "DNT": "1",
        "Origin": "https://xt.bjwxdxh.org.cn",
        "Referer": "https://xt.bjwxdxh.org.cn/static/member/",
        "Sec-Fetch-Dest": "empty",
        "Sec-Fetch-Mode": "cors",
        "Sec-Fetch-Site": "same-origin",
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
        "sec-ch-ua": '"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"',
        "sec-ch-ua-mobile": "?0",
        "sec-ch-ua-platform": '"Linux"',
    }
    data = {
        "username": f"{username}",
        "password": f"{password}",
        "type": "account",
    }
    response = session.post(
        "https://xt.bjwxdxh.org.cn/loginapi/login", headers=headers, data=data
    )
    if os.getenv("DEBUG"):
        print(response.text)
    if response.status_code != 200:
        print("Fail to login")
        exit(1)
    if response.json()["ret"] != 1:
        print("Fail to login, check your username/password")
        exit(1)
    print("Login succeed!")


def get_possible_plan(session):
    print("Try to select exam plan...")
    headers = {
        "Accept": "*/*",
        "Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
        "Connection": "keep-alive",
        "DNT": "1",
        "Referer": "https://xt.bjwxdxh.org.cn/static/member/",
        "Sec-Fetch-Dest": "empty",
        "Sec-Fetch-Mode": "cors",
        "Sec-Fetch-Site": "same-origin",
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
        "sec-ch-ua": '"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"',
        "sec-ch-ua-mobile": "?0",
        "sec-ch-ua-platform": '"Linux"',
    }

    params = {
        "pageIndex": "1",
    }

    response = session.get(
        "https://xt.bjwxdxh.org.cn/memberapi/getExamPlanCanSelected",
        params=params,
        headers=headers,
    )
    if os.getenv("DEBUG"):
        print(response.text)
    if response.status_code != 200:
        print("Fail to select exam plan")
        exit(1)
    result = response.json()
    possible_plans = result["data"]["list"]
    plan_id = possible_plans[0]["id"]
    print("Selected:", possible_plans[0]["title"])
    return plan_id


def reserve(session, plan_id):
    print("Try to reserve the exam...")
    headers = {
        "Accept": "application/json",
        "Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
        "Connection": "keep-alive",
        "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
        "DNT": "1",
        "Origin": "https://xt.bjwxdxh.org.cn",
        "Referer": "https://xt.bjwxdxh.org.cn/static/member/",
        "Sec-Fetch-Dest": "empty",
        "Sec-Fetch-Mode": "cors",
        "Sec-Fetch-Site": "same-origin",
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
        "sec-ch-ua": '"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"',
        "sec-ch-ua-mobile": "?0",
        "sec-ch-ua-platform": '"Linux"',
    }

    data = {
        "id": f"{plan_id}",
    }
    response = session.post(
        "https://xt.bjwxdxh.org.cn/memberapi/addMemberExamPlan",
        headers=headers,
        data=data,
    )
    if os.getenv("DEBUG"):
        print(response.text)
    if response != 200:
        print("Fail to reserve the exam!")
        exit(1)
    print("OK!")


args = sys.argv

if len(args) != 3:
    print("Usage: ./prog <username> <password>")
with requests.Session() as session:
    login(session, args[1], args[2])
    if os.getenv("DEBUG"):
        for cookie in session.cookies:
            print(cookie)
    plan_id = get_possible_plan(session)
    reserve(session, plan_id)
