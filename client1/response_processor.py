import codecs

def keyword_processor(keyword):
    newkeyword = ""

    for i in keyword:
        b = str(codecs.encode(bytes(i,"utf-8"),"hex"))[2:4]
        if b == "27":
            continue
        elif b == '22' or b == '23' or b == '24' or b == '25' or b == '36' or b == '2b' or b == '2c' or b == '2f' or b == '3a' or b == '3c' or b == '3b' or b == '3d' or b == '3e' or b == '3f' or b == '40':
            a = "%" + b
            newkeyword += a
        elif b == '5b' or b == '5c' or b == '5d' or b == '5e' or b == '60' :
            a = "%"+b
            newkeyword += a
        else:
            newkeyword += i

    return newkeyword
    #print(results)