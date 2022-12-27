#! /usr/bin/env python3
import itertools

digitvalues = {
    '2': 2,
    '1': 1,
    '0': 0,
    '-': -1,
    '=': -2,
}

inversevalues = dict((v, k) for k, v in digitvalues)

def fromSnafu(s):
    placevalue = 1
    n = 0
    while len(s) > 0:
        digit = s[-1]
        s=s[:-1]
        n += digitvalues[digit]*placevalue
        placevalue *= 5
    return n

# this routine isn't finished yet
def toSnafu(n):
    placevalue = 1
    s = ""
    while n >= 3*placevalue:
        placevalue *= 5
    while True:
        if n in inversevalues:
            s += inversevalues[n]
            break
        for k, v in inversevalues:
            diff = k*placevalue+k*placevalue/5
            if diff <= n:
                s += v
                n -= placevalue
    return s



def part1(data):
    return sum([fromSnafu(d) for d in data])

if __name__ == "__main__":
    f = open("./inputsample.txt")
    lines = f.readlines()
    data = [l.strip() for l in lines]
    print(part1(data))