#! /usr/bin/env python3
import itertools
import collections

def cube(l):
    return [float(x) for x in l.split(',')]

def faces(cube):
    return [
        ('x', cube[0]-0.5, cube[1], cube[2]),
        ('x', cube[0]+0.5, cube[1], cube[2]),
        ('y', cube[0], cube[1]-0.5, cube[2]),
        ('y', cube[0], cube[1]+0.5, cube[2]),
        ('z', cube[0], cube[1], cube[2]-0.5),
        ('z', cube[0], cube[1], cube[2]+0.5),
    ]

openfaces=set()
def collectFace(f):
    if f in openfaces:
        openfaces.remove(f)
    else:
        openfaces.add(f)

def minmax(r):
    return min(r), max(r)

if __name__ == "__main__":
    f = open("./input.txt")
    lines = f.readlines()
    cubes = [cube(l.strip()) for l in lines]
    for c in cubes:
        for f in faces(c):
            collectFace(f)
    count = len(openfaces)
    print(count)
