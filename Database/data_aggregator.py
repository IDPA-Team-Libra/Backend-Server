

def read_file(file):
    f = open(file, "r")
    return f.readline()


symbol = read_file("symbol.txt").split(",")
names = read_file("company_names.txt").split(",")
print(symbol)
for x in range(0, len(symbol)):
    form = "('"+symbol[x]+"','"+names[x] + "', '5', NOW()),"
    print(form)
