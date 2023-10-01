import json


with open('new_products.json') as fi:
	file_content = fi.read()

parsed = json.loads(file_content)

print(parsed[0]['manufacturer'])

cleaned = []

for item in parsed:
	simple_item = {}
	simple_item['sku'] = item['sku']
	simple_item['name'] = item['name']
	simple_item['type'] = item['type']
	simple_item['price'] = item['price']
	simple_item['category'] = item['category']
	simple_item['manufacturer'] = item['manufacturer']
	cleaned.append(simple_item)

with open('cleaned_products.json', 'w') as outfile:
	json.dump(parsed, outfile)


print('done')
