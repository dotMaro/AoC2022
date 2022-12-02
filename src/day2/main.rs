use std::fs;

fn main() {
    let input =
        fs::read_to_string("src/day2/input.txt").expect("Should be able to read input file");
    println!("Part 1. The score would be {}", total_score(input.as_str()));
    println!(
        "Part 2. The score would be {}",
        total_score_when_result_is_fixed(input.as_str())
    );
}

fn total_score(input: &str) -> u64 {
    input
        .lines()
        .map(|line| {
            let mut chars = line.chars();
            let opponent_sign_code = chars.next().unwrap();
            let your_sign_code = chars.last().unwrap();
            Sign::new(your_sign_code).score(&Sign::new(opponent_sign_code))
        })
        .sum()
}

fn total_score_when_result_is_fixed(input: &str) -> u64 {
    input
        .lines()
        .map(|line| {
            let mut chars = line.chars();
            let opponent_sign_code = chars.next().unwrap();
            let result_code = chars.last().unwrap();
            let opponent_sign = Sign::new(opponent_sign_code);
            Sign::new_to_get_result_against(Result::new(result_code), &opponent_sign)
                .score(&opponent_sign)
        })
        .sum()
}

enum Sign {
    Rock,
    Paper,
    Scissors,
}

enum Result {
    Win,
    Loss,
    Draw,
}

impl Result {
    fn new(c: char) -> Result {
        match c {
            'X' => Result::Loss,
            'Y' => Result::Draw,
            'Z' => Result::Win,
            _ => panic!("invalid result char {}", c),
        }
    }
}

impl Sign {
    fn new(c: char) -> Sign {
        match c {
            'A' | 'X' => Sign::Rock,
            'B' | 'Y' => Sign::Paper,
            'C' | 'Z' => Sign::Scissors,
            _ => panic!("invalid sign char {}", c),
        }
    }

    fn new_to_get_result_against(res: Result, other: &Sign) -> Sign {
        match res {
            Result::Draw => match other {
                Sign::Rock => Sign::Rock,
                Sign::Paper => Sign::Paper,
                Sign::Scissors => Sign::Scissors,
            },
            Result::Win => match other {
                Sign::Rock => Sign::Paper,
                Sign::Paper => Sign::Scissors,
                Sign::Scissors => Sign::Rock,
            },
            Result::Loss => match other {
                Sign::Rock => Sign::Scissors,
                Sign::Paper => Sign::Rock,
                Sign::Scissors => Sign::Paper,
            },
        }
    }

    fn score(&self, other: &Sign) -> u64 {
        let shape_score = match self {
            Sign::Rock => 1,
            Sign::Paper => 2,
            Sign::Scissors => 3,
        };

        let outcome_score = match (self, other) {
            // Draw.
            (Sign::Rock, Sign::Rock) => 3,
            (Sign::Paper, Sign::Paper) => 3,
            (Sign::Scissors, Sign::Scissors) => 3,
            // Win.
            (Sign::Rock, Sign::Scissors) => 6,
            (Sign::Paper, Sign::Rock) => 6,
            (Sign::Scissors, Sign::Paper) => 6,
            // Loss.
            (Sign::Rock, Sign::Paper) => 0,
            (Sign::Paper, Sign::Scissors) => 0,
            (Sign::Scissors, Sign::Rock) => 0,
        };

        shape_score + outcome_score
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn total_score_with_example() {
        let example = "A Y\nB X\nC Z";
        assert!(total_score(example) == 15);
    }
}
