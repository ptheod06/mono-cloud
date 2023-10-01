import json
from datetime import datetime
import math


with open('cleaned_products.json') as fi:
        file_content = fi.read()

parsed = json.loads(file_content)

common_cats = []

starttime = datetime.now()

for i in range(0, 6000):
	inner_arr = []
	for j in range(0, 6000):

		if i == j:
			inner_arr.append('na')
			continue
		else:
			set1 = set(parsed[i]['category'])
			set2 = set(parsed[j]['category'])
			common = set1.intersection(set2)
			inter_list = list(common)
			similarity = (len(inter_list) / (math.sqrt(len(parsed[i]['category'])) + math.sqrt(len(parsed[j]['category']))))
			inner_arr.append(similarity)

	common_cats.append(inner_arr)

endtime = datetime.now()

with open('common_categories.json', 'w') as fo:
	json.dump(common_cats, fo)

print(endtime - starttime)
print('done')
