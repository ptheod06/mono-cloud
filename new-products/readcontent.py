import json
from datetime import datetime
import math


with open('cleaned_products.json') as fi:
        file_content = fi.read()

parsed = json.loads(file_content)


print(parsed[0])
