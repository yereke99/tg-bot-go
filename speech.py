from speechkit import SpeechSynthesis, Session
import requests

# Настройка параметров
IAM_TOKEN = 't1.9euelZrIm8_ImMzOxpXJi47HjZzOne3rnpWayovIkZbLk8zOic2byZSVx8jl9PdSRzFD-e8TKgHO3fT3EnYuQ_nvEyoBzs3n9euelZqcm4uRjpiejYqPzZPLyJONne_8xeuelZqcm4uRjpiejYqPzZPLyJONnQ.2irZ28qhBPb15rTZIqVOnWnaf9wv78NVHRwnQNTigL7N-Fu9R9uEzKuOkal1XV8rYs2oYXS5hJ5K4gui3k74CQ'
FOLDER_ID = 'b1gk6dof1frm98hckkch'


url = 'https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize'
headers = {
    'Authorization': f'Bearer {IAM_TOKEN}'
}
data = {
    'text': 'Салам Богдан! Как дела? Что делаешь дрочишь или нет?',
    'lang': 'kk-KK',  # Казахский язык
    'voice': 'amira',  # Голос для казахского языка
    'folderId': FOLDER_ID
}

datas = {
    'text': 'Я могу конвертировать казахский текст на казахский аудио файл, тем более это популярный и ценный функциянал в Казахстане',
    'lang': 'ru-RU',  # Руский язык
    'voice': 'alena',  # Голос для казахского языка
    'folderId': FOLDER_ID
}

response = requests.post(url, headers=headers, data=datas)
if response.status_code == 200:
    with open('result.ogg', 'wb') as f:
        f.write(response.content)
    print("Аудио успешно сохранено как result.ogg")
else:
    print(f"Ошибка: {response.status_code}, {response.text}")