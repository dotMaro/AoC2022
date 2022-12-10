use std::{collections::HashSet, fs};

fn main() {
    let input = fs::read_to_string("src/day9/input.txt").unwrap();
    println!(
        "Part 1. The tail has been in {} positions.",
        execute_instructions(&input, 2)
    );
    println!(
        "Part 2. The tail has been in {} positions.",
        execute_instructions(&input, 10)
    );
}

fn execute_instructions(input: &str, knots: usize) -> usize {
    let mut rope = Rope::new(knots);
    input
        .lines()
        .map(|line| {
            let words = line.split_once(' ').unwrap();
            (
                Direction::new(words.0.chars().next().unwrap()),
                words.1.parse().unwrap(),
            )
        })
        .for_each(|(d, n)| rope.step_n(&d, n));
    rope.tail_positions_visited.len()
}

struct Rope {
    knots: Vec<Coord>,
    tail_positions_visited: HashSet<Coord>,
}

#[derive(Debug)]
enum Direction {
    Right,
    Left,
    Up,
    Down,
}

impl Direction {
    fn new(c: char) -> Self {
        match c {
            'R' => Direction::Right,
            'L' => Direction::Left,
            'U' => Direction::Up,
            'D' => Direction::Down,
            _ => panic!(),
        }
    }
}

impl Rope {
    fn new(knots: usize) -> Self {
        let mut tail_positions_visited = HashSet::new();
        tail_positions_visited.insert(Coord { x: 0, y: 0 });
        Rope {
            knots: vec![Coord { x: 0, y: 0 }; knots],
            tail_positions_visited,
        }
    }

    fn step_n(&mut self, direction: &Direction, n: u8) {
        for _ in 0..n {
            self.step(direction);
        }
    }

    fn step(&mut self, direction: &Direction) {
        match direction {
            Direction::Right => self.knots[0].x += 1,
            Direction::Left => self.knots[0].x -= 1,
            Direction::Up => self.knots[0].y += 1,
            Direction::Down => self.knots[0].y -= 1,
        }

        for i in 1..self.knots.len() {
            if self.knots[i].distance_larger_than_one(&self.knots[i - 1]) {
                self.adjust_tail(i);
            }
        }
    }

    fn adjust_tail(&mut self, i: usize) {
        let target = self.knots[i - 1].clone();
        let knot_len = self.knots.len();
        let knot = &mut self.knots[i];
        if target.x > knot.x {
            knot.x += 1;
        } else if target.x < knot.x {
            knot.x -= 1;
        }
        if target.y > knot.y {
            knot.y += 1;
        } else if target.y < knot.y {
            knot.y -= 1;
        }
        if i == knot_len - 1 {
            self.tail_positions_visited.insert(Coord {
                x: knot.x,
                y: knot.y,
            });
        }
    }
}

#[derive(PartialEq, Eq, Hash, Clone)]
struct Coord {
    x: i64,
    y: i64,
}

impl Coord {
    fn distance_larger_than_one(&self, other: &Coord) -> bool {
        self.x.abs_diff(other.x) > 1 || self.y.abs_diff(other.y) > 1
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn examples_part1() {
        let example = "R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2";
        let tail_pos_count = execute_instructions(example, 2);
        assert_eq!(tail_pos_count, 13);
    }

    #[test]
    fn examples_part2() {
        let example = "R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2";
        let tail_pos_count = execute_instructions(example, 10);
        assert_eq!(tail_pos_count, 1);
    }
}
