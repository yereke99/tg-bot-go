import socket
import random
from ipaddress import ip_network
from concurrent.futures import ThreadPoolExecutor

def get_random_ip(network):
    """Generates a random IP address from the given network range."""
    return str(random.choice(list(ip_network(network).hosts())))

def scan_port(ip, port):
    """Scans a specific port on the given IP address."""
    try:
        with socket.create_connection((ip, port), timeout=1):
            return True
    except (socket.timeout, ConnectionRefusedError, OSError):
        return False

def attempt_backdoor(ip, port):
    """Placeholder function for attempting to exploit a backdoor (for educational purposes only)."""
    # WARNING: This is for educational purposes only. Do not use this for illegal activities.
    print(f"[*] Attempting to exploit backdoor on {ip}:{port}...")
    # Simulate a backdoor attempt (replace with actual code if authorized)
    # Example: Try to connect and send a payload
    try:
        with socket.create_connection((ip, port), timeout=1) as sock:
            sock.sendall(b"Malicious payload\n")
            response = sock.recv(1024)
            if response:
                print(f"[+] Backdoor successful on {ip}:{port}")
                return True
    except Exception as e:
        print(f"[-] Backdoor failed on {ip}:{port}: {e}")
    return False

def scan_random_server(networks, ports):
    """Scans a random IP address and port from the given networks and ports."""
    network = random.choice(networks)
    ip = get_random_ip(network)
    port = random.choice(ports)
    if scan_port(ip, port):
        print(f"[+] Open server found: {ip}:{port}")
        return ip, port
    return None

def find_open_servers(networks, ports, thread_count=20, scan_attempts=100):
    """Finds open servers using multithreading and a limited number of scan attempts."""
    open_servers = []
    
    with ThreadPoolExecutor(max_workers=thread_count) as executor:
        futures = [executor.submit(scan_random_server, networks, ports) for _ in range(scan_attempts)]
        for future in futures:
            result = future.result()
            if result:
                open_servers.append(result)
                # Attempt to exploit backdoor (for educational purposes only)
                attempt_backdoor(result[0], result[1])

    # Print only open servers
    if open_servers:
        print("\n=== Found Open Servers ===")
        for ip, port in open_servers:
            print(f"{ip}:{port}")
    else:
        print("\n[!] No open servers found.")

if __name__ == "__main__":
    if __name__ == "__main__":
        # Define IP ranges for Kazakhstan and other regions
        networks = [
            # Kazakhstan
            "5.34.160.0/19",
            "95.56.0.0/14",
            "2.132.0.0/14",
            "89.218.0.0/16",
            "195.12.240.0/20",
            "212.154.160.0/20",
            
            # Additional Global IP Ranges (for educational purposes only)
            "1.0.0.0/24",        # Example range
            "8.8.8.0/24",        # Google Public DNS (example)
            "104.16.0.0/12",     # Cloudflare (example)
            "172.217.0.0/16",    # Google (example)
            "192.168.0.0/16",    # Private network range (example)
            "10.0.0.0/8",        # Private network range (example)
            "203.0.113.0/24",    # Reserved for documentation (example)
            "198.51.100.0/24",   # Reserved for documentation (example)
            "93.184.216.0/24",   # Example range
            "45.60.0.0/16",      # Example range
            "74.125.0.0/16",     # Google (example)
            "216.58.0.0/16",     # Google (example)
            "64.233.160.0/19",   # Google (example)
            "66.102.0.0/20",     # Google (example)
            "72.14.192.0/18",    # Google (example)
            "209.85.128.0/17",   # Google (example)
            "173.194.0.0/16",    # Google (example)
            "207.126.144.0/20",  # Microsoft (example)
            "40.112.0.0/13",     # Microsoft (example)
            "52.96.0.0/12",      # Microsoft Azure (example)
            "104.146.0.0/15",    # Microsoft (example)
            "131.253.0.0/16",    # Microsoft (example)
            "137.116.0.0/15",    # Microsoft (example)
            "168.63.0.0/16",     # Microsoft (example)
            "13.64.0.0/11",      # Microsoft Azure (example)
            "40.74.0.0/15",      # Microsoft (example)
            "52.224.0.0/12",     # Microsoft (example)
            "65.52.0.0/14",      # Microsoft (example)
            "94.245.0.0/16",     # Example range
            "185.94.0.0/16",     # Example range
            "212.83.0.0/16",     # Example range
            "217.69.0.0/16",     # Example range
        ]

        # Define target ports to scan
        target_ports = [22, 80, 443, 8080, 3306, 21, 23, 25, 110, 143, 3389]

        # Set thread count and scan attempts
        thread_count = 20
        scan_attempts = 200  # Number of IP addresses to scan

        print("Starting open server search...\n")
        find_open_servers(networks, target_ports, thread_count, scan_attempts)