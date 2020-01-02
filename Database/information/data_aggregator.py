def read_file(file):
    f = open(file, "r")
    return f.readline()


symbol = read_file("symbol.txt").split(",")
names = read_file("company_names.txt").split(",")
print(symbol)
for x in range(0, len(symbol)):
    ranges = ["1", "5", "15", "30", "60"]
    for y in ranges:
        form = "('"+symbol[x]+"','{}', NOW()),".format(y)
        print(form)
