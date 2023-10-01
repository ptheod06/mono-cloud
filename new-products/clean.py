import json


with open('new_products.json') as fi:
        file_content = fi.read()

parsed = json.loads(file_content)

for item in parsed:
	del item['url']
	del item['image']

with open('cleaned_products.json', 'w') as fo:
	json.dump(parsed, fo)

print('done')
