#! /usr/bin/env python3
import itertools

# this order is important
digitvalues = {
    '2': 2,
    '=': -2,
    '1': 1,
    '-': -1,
    '0': 0,
}

# inversevalues = dict((v, k) for k, v in digitvalues.items())

def fromSnafu(s):
    placevalue = 1
    n = 0
    while len(s) > 0:
        digit = s[-1]
        s=s[:-1]
        n += digitvalues[digit]*placevalue
        placevalue *= 5
    return n

def sign(n):
    if n > 0:
        return 1
    elif n < 0:
        return -1
    return 0


# this routine isn't finished yet
def toSnafu(n):
    # for simple cases just get out early
    # if n in inversevalues:
    #     return inversevalues[n]

    # 222 is 50+10+2 == 62 which is one less than half of 125
    # 1=== is 125-50-10-2 == 63 which is one more than half of 125
    # we know that positive numbers, which is all we have, *always* need a first
    # digit of 1 or 2, so let's start with figuring that out
    placevalue = 1
    while n >= placevalue*5//2:
        placevalue *= 5

    s = ""
    while placevalue > 0:
        incoming = n
        if n > 0:
            if n > placevalue + placevalue//2:
                s += "2"
                n -= placevalue*2
            elif n > placevalue//2:
                s += "1"
                n -= placevalue
            else:
                s += "0"
        elif n < 0:
            if n < -placevalue - placevalue//2:
                s += "="
                n += placevalue*2
            elif n < -placevalue//2:
                s += "-"
                n += placevalue
            else:
                s += "0"
        else: # n == 0:
            s += "0"
        # print(f"incoming={incoming}, now={n}, pv={placevalue}, s='{s}'")
        placevalue //= 5

    return s



def part1(data):
    # for d in data:
    #     n = fromSnafu(d)
    #     s = toSnafu(n)
    #     print(f"> original: {d} decimal: {n}  back: {s}")

    return sum([fromSnafu(d) for d in data])

if __name__ == "__main__":
    f = open("./input.txt")
    lines = f.readlines()
    data = [l.strip() for l in lines]
    v = part1(data)
    print(v, toSnafu(v))