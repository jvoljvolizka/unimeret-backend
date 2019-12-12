import time, subprocess

subprocess.Popen([ "env" ,"GOOS=linux", "GOARCH=amd64" ,"go" ,"build"])

time.sleep(3)

st = subprocess.check_output((["ls"]))
st = st.decode()

st = st.split("\n")

st.pop(-1)
out = ""

for i in st:
    if ".go" not in i:
        out = i

subprocess.Popen([ "zip" , "-j" , "./unimeret-" + out ,  "./" + out ])

time.sleep(1)

loc = subprocess.check_output((["pwd"]))

loc = loc.decode()
loc = loc.strip("\n")
test = subprocess.check_output([ "aws" , "lambda" , "update-function-code" ,  "--function-name" , "unimeret-" +  out  , "--zip-file" , "fileb://" + loc + "/unimeret-" +  out + ".zip"])

test = test.decode()

print(test)


subprocess.Popen([ "rm" ,  "unimeret-"  + out + ".zip" ])
