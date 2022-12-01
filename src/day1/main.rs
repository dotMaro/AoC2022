use std::fs;

fn main() {
    let input = fs::read_to_string("src/day1/input.txt").expect("Should be able to read the input");
    let sorted_calories = sorted_calories(input.as_str());
    println!(
        "Part 1. The elf with the most calories has {} calories",
        most_calories(&sorted_calories)
    );
    println!(
        "Part 2. The sum of the top 3 elves with the most calories is {}",
        sum_of_top_n_calories(3, &sorted_calories)
    );
}

fn most_calories(sorted_calories: &Vec<u64>) -> u64 {
    *sorted_calories.last().unwrap()
}

fn sum_of_top_n_calories(n: usize, sorted_calories: &Vec<u64>) -> u64 {
    sorted_calories[sorted_calories.len() - n..]
        .into_iter()
        .sum()
}

fn sorted_calories(input: &str) -> Vec<u64> {
    let mut calories = vec![];
    let mut current: u64 = 0;
    for line in input.lines() {
        if line == "" {
            calories.push(current);
            current = 0;
            continue;
        }

        current += line.parse::<u64>().unwrap();
    }
    calories.push(current);

    calories.sort();
    calories
}
