import json
from datetime import datetime
import math


with open('common_categories.json') as fi:
        file_content = fi.read()

parsed = json.loads(file_content)

for i in range(0, 200):
	print(i, parsed[123][i])

#print(parsed[0][2])
#print(parsed[2][3])
