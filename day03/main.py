#! /usr/bin/env python3
import itertools

priorities = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"


def perElf(d):
    l = len(d) // 2
    # print(d, l)
    d1 = set(d[:l])
    d2 = set(d[l:])
    common = d1.intersection(d2)
    return priorities.index("".join([c for c in common])) + 1


def perGroup(g):
    common = set(g[0]).intersection(set(g[1])).intersection(set(g[2]))
    return priorities.index("".join([c for c in common])) + 1


def part1(data):
    s = sum([perElf(d) for d in data])
    return s


def part2(data):
    groups = [
        (data[i * 3], data[i * 3 + 1], data[i * 3 + 2]) for i in range(len(data) // 3)
    ]
    s = sum([perGroup(g) for g in groups])
    return s


if __name__ == "__main__":
    f = open("./inputsample.txt")
    lines = f.readlines()
    data = [l.strip() for l in lines]
    print(part1(data))
    print(part2(data))
