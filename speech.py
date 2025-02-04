from speechkit import SpeechSynthesis, Session
import requests

# Настройка параметров
IAM_TOKEN = 't1.9euelZqeipqSmIuQipaNic-PlM3Miu3rnpWayovIkZbLk8zOic2byZSVx8jl8_d5OhdD-e8adWUM_d3z9zlpFEP57xp1ZQz9zef1656VmsnNnsvKjsbMms2Kyo_Nz5jO7_zF656VmsnNnsvKjsbMms2Kyo_Nz5jO.kgcbngC9mgOzMOxy0cX_Feb-lzuJz7fuUR2xGdOnKzaKGSUJES1YUzRbXvEmNlVsNZ8HlC4fnSb51KI_B1foAw'
FOLDER_ID = 'b1gk6dof1frm98hckkch'


url = 'https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize'
headers = {
    'Authorization': f'Bearer {IAM_TOKEN}'
}
data = {
    'text': 'Қайырлы күн менің есімім Амира! Мен жасанды интелектпін. Қазақша сөйлей аламын, бар болғаны маған керекті мәтінді жазсаңыз болғаны)',
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

response = requests.post(url, headers=headers, data=data)
if response.status_code == 200:
    with open('result.ogg', 'wb') as f:
        f.write(response.content)
    print("Аудио успешно сохранено как result.ogg")
else:
    print(f"Ошибка: {response.status_code}, {response.text}")