import json
from datetime import datetime
import math


with open('last_products.json') as fi:
        file_content = fi.read()

parsed = json.loads(file_content)


for item in parsed:
	if item['price'] <= 40:
		item['price'] = 0
	elif item['price'] <= 80:
		item['price'] = 1
	elif item['price'] <= 120:
		item['price'] = 2
	elif item['price'] <= 160:
		item['price'] = 3
	elif item['price'] <= 200:
		item['price'] = 4
	elif item['price'] <= 240:
		item['price'] = 5
	elif item['price'] <= 280:
		item['price'] = 6
	else:
		item['price'] = 7


with open('final_products.json', 'w') as fo:
	json.dump(parsed, fo)
