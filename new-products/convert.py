import json


with open('fixmanu_products.json') as fi:
	file_content = fi.read()

parsed = json.loads(file_content)


cleaned = []

for item in parsed:
	simple_item = {}
	simple_item['sku'] = item['sku']
	simple_item['name'] = item['name']
	simple_item['type'] = item['type']
	simple_item['price'] = item['price']
	simple_item['manufacturer'] = item['manufacturer']

	category_names = []

	for category in item['category']:
		category_names.append(category['name'])

	simple_item['category'] = category_names
	cleaned.append(simple_item)

with open('cleaned_products.json', 'w') as outfile:
	json.dump(cleaned, outfile)


print('done')
