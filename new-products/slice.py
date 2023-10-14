import json

with open('final_products.json') as fi:
        file_content = fi.read()



parsed = json.loads(file_content)


cleaned = []

cleaned.append(parsed[0])
cleaned.append(parsed[5])
cleaned.append(parsed[25])


with open('first_products.json', 'w') as outfile:
        json.dump(cleaned, outfile)
