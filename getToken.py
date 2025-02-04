#!/usr/bin/env python3

import requests

def main():
    # Подставьте свой OAuth-токен:
    token = "y0__xDhxNm2BBjB3RMgjarxjBKNtEqZiQE-K-a_Ae-YVJ10EKvabg"

    url = "https://iam.api.cloud.yandex.net/iam/v1/tokens"
    payload = {
        "yandexPassportOauthToken": token
    }

    try:
        response = requests.post(url, json=payload)
        response.raise_for_status()  # выбросит исключение, если код ответа != 200..299

        # Если всё в порядке, печатаем полученный JSON (IAM-токен и дополнительную информацию)
        print("Успешный ответ:")
        print(response.json())

    except requests.exceptions.RequestException as e:
        print("Произошла ошибка при запросе:")
        print(e)

if __name__ == "__main__":
    main()
