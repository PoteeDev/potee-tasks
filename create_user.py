import requests
from dataclasses import dataclass
import sys
url = "http://localhost"
url = "http://192.168.0.20"
url = sys.argv[1]
@dataclass
class User:
    first_name: str
    second_name: str
    group: str
    password: str
    email: str

    def register(self):
        r = requests.post(f"{url}/register", json=self.__dict__)
        print("code:", r.status_code, "answer:", r.text)

def test_user():
    user = User("Ivan", "Ivanov", "BBSO-03-17", "secretpass", "iivanov@gmail.com")
    user.register()

if __name__ == "__main__":
    test_user()