import json


with open('products.json') as fi:
	file_content = fi.read()

parsed = json.loads(file_content)

count = 0

for item in parsed:
	if 'manufacturer' not in item:
		item['manufacturer'] = 'unknown'


with open('fixmanu_products.json', 'w') as fo:
	json.dump(parsed, fo)

print('done')
