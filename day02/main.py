#! /usr/bin/env python3
import itertools

# The score for a single round is the score for the shape
# you selected (1 for Rock, 2 for Paper, and 3 for Scissors)
# plus the score for the outcome of the round (0 if you lost,
# 3 if the round was a draw, and 6 if you won).
scores = {
    "AA": 1 + 3,
    "AB": 2 + 6,
    "AC": 3 + 0,
    "BA": 1 + 0,
    "BB": 2 + 3,
    "BC": 3 + 6,
    "CA": 1 + 6,
    "CB": 2 + 0,
    "CC": 3 + 3,
}


def part1(data):
    mapping = {"X": "A", "Y": "B", "Z": "C"}
    d2 = [d[0] + mapping[d[1]] for d in data]
    s = sum([scores[d] for d in d2])
    print(s)


def part2(data):
    move = {
        "AX": "C",
        "AY": "A",
        "AZ": "B",
        "BX": "A",
        "BY": "B",
        "BZ": "C",
        "CX": "B",
        "CY": "C",
        "CZ": "A",
    }
    d2 = [d[0] + move[d] for d in data]
    s = sum([scores[d] for d in d2])
    print(s)


if __name__ == "__main__":
    f = open("./input.txt")
    lines = f.readlines()
    data = [l[0] + l[2] for l in lines]
    part1(data)
    part2(data)
