#! /usr/bin/env python3
import itertools
import collections

def cube(tup):
    return [float(x) for x in tup]

# for each cube, rather than worrying about neighbors, we can define faces that
# live between the cubes on each axis and will be positioned in a way that two
# touching faces have equal values that can serve as keys in a set or a map.
def faces(cube):
    return [
        ('x', cube[0]-0.5, cube[1], cube[2]),
        ('x', cube[0]+0.5, cube[1], cube[2]),
        ('y', cube[0], cube[1]-0.5, cube[2]),
        ('y', cube[0], cube[1]+0.5, cube[2]),
        ('z', cube[0], cube[1], cube[2]-0.5),
        ('z', cube[0], cube[1], cube[2]+0.5),
    ]

def collectFace(f, openset):
    if f in openset:
        openset.remove(f)
    else:
        openset.add(f)

def minmax(r):
    return min(r), max(r)

def part1(cubes):
    '''
    Collect all the faces into a set, except that when we go to add one, if
    it's already there we delete it.
    '''
    openfaces=set()

    for c in cubes:
        for f in faces(c):
            collectFace(f, openfaces)
    count = len(openfaces)
    return count

def part2(cubes, opencount):
    '''
    Now that we have the set of cubes, we can fill in all the empty cubes in its vicinity
    with fake cubes, and then iterate multiple times, removing all of the fake cubes that
    have faces with no (real or fake) neighbors, which are the ones that touch outside air.
    Once we can't find any more, the remaining fake cubes must be interior cubes, so we can
    count the exposed faces on them and subtract that from the face count.
    '''
    minx, maxx = minmax([c[0] for c in cubes])
    miny, maxy = minmax([c[1] for c in cubes])
    minz, maxz = minmax([c[2] for c in cubes])

    fakecubes = []
    for x in range(int(minx-1), int(maxx+2)):
        for y in range(int(miny-1), int(maxy+2)):
            for z in range(int(minz-1), int(maxz+2)):
                c = cube((x, y, z))
                if c not in cubes:
                    fakecubes.append(c)

    while True:
        fakefaces=collections.defaultdict(int)
        for c in fakecubes:
            for f in faces(c):
                fakefaces[f]+=1
        for c in cubes:
            for f in faces(c):
                fakefaces[f] += 1

        removes = []
        for c in fakecubes:
            for f in faces(c):
                if fakefaces[f] == 1:
                    removes.append(c)

        if len(removes) == 0:
            break

        fakecubes = [c for c in fakecubes if c not in removes]

    # print(fakecubes)
    fakeopenfaces=set()
    for c in fakecubes:
        for f in faces(c):
            collectFace(f, fakeopenfaces)

    return opencount - len(fakeopenfaces)

if __name__ == "__main__":
    f = open("./input.txt")
    lines = f.readlines()
    cubes = [cube(l.strip().split(',')) for l in lines]
    opencount = part1(cubes)
    print(opencount)
    print(part2(cubes, opencount))