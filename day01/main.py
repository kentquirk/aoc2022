#! /usr/bin/env python3
import itertools


def parse(lines):
    data = [l.strip() for l in lines]
    data.append("")
    # print(f"{data}")

    elves = []
    f = []
    for d in data:
        if d == "" and len(f) > 0:
            elves.append([int(s) for s in f])
            f = []
        else:
            f.append(d)

    return elves


def part1(lines):
    elves = parse(lines)
    # print(f"{elves}")

    best = 0
    for e in elves:
        s = sum(e)
        if s > best:
            best = s

    print(best)


def part2(lines):
    elves = parse(lines)
    # print(f"{elves}")

    elvesSort = sorted(elves, key=lambda e: sum(e), reverse=True)
    top3 = elvesSort[:3]
    s = sum([sum(e) for e in top3])
    print(s)


if __name__ == "__main__":
    f = open("./input.txt")
    lines = f.readlines()
    part1(lines)
    part2(lines)
