import socket
import random
from ipaddress import ip_network
from concurrent.futures import ThreadPoolExecutor

def get_random_ip(network):
    """Генерирует случайный IP-адрес из заданного диапазона."""
    return str(random.choice(list(ip_network(network).hosts())))

def scan_port(ip, port):
    """Проверяет, открыт ли порт на указанном IP."""
    try:
        with socket.create_connection((ip, port), timeout=1):
            return True
    except (socket.timeout, ConnectionRefusedError):
        return False

def scan_random_server(networks, ports):
    """Сканирует случайный IP-адрес и порт."""
    network = random.choice(networks)
    ip = get_random_ip(network)
    port = random.choice(ports)
    if scan_port(ip, port):
        print(f"[+] Найден открытый сервер: {ip}:{port}")
        return ip, port
    else:
        print(f"[-] {ip}:{port} закрыт")
        return None

def find_open_servers(networks, ports, thread_count=10):
    """Ищет открытые сервера с многопоточностью."""
    with ThreadPoolExecutor(max_workers=thread_count) as executor:
        while True:
            executor.submit(scan_random_server, networks, ports)

if __name__ == "__main__":
    # Расширенный список диапазонов IP
    kazakhstan_networks = [
        "5.34.160.0/19",
        "95.56.0.0/14",
        "2.132.0.0/14",
        "89.218.0.0/16",
        "195.12.240.0/20",
        "212.154.160.0/20",
    ]
    # Порты для проверки
    target_ports = [22, 80, 443, 8080, 3306, 21, 23, 25, 110, 143, 3389]  # Расширенные порты

    # Количество потоков
    thread_count = 20

    print("Запускаем поиск открытых серверов...")
    find_open_servers(kazakhstan_networks, target_ports, thread_count)
