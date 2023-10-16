import json
from datetime import datetime
import math


with open('final_products.json') as fi:
        file_content = fi.read()

parsed = json.loads(file_content)

common_cats = []

starttime = datetime.now()

for i in range(0, 15):
	inner_arr = []
	for j in range(0, 15):

		if i == j:
			inner_arr.append(-1)
			continue
		else:
			set1 = set(parsed[i]['category'])
			set2 = set(parsed[j]['category'])
			common = set1.intersection(set2)
			inter_list = list(common)
			denominator = math.sqrt(len(parsed[i]['category']) + 3) * math.sqrt(len(parsed[j]['category']) + 3)
			numerator = len(inter_list)
			if parsed[i]['price'] == parsed[j]['price']:
				numerator += 1

			if parsed[i]['type'] == parsed[j]['type']:
				numerator += 1
			try:
				if parsed[i]['manufacturer'] == parsed[j]['manufacturer']:
					numerator += 1
			except:
				print(parsed[i])
				print(parsed[j])
				exit()
			similarity = numerator / denominator
			inner_arr.append(similarity)

	common_cats.append(inner_arr)

endtime = datetime.now()

with open('common_categories.json', 'w') as fo:
	json.dump(common_cats, fo, indent=2)

print(endtime - starttime)
print('done')
