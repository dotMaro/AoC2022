use std::fs;

fn main() {
    let input = fs::read_to_string("src/day5/input.txt").expect("Should be able to read the input");
    let stacks = parse_stacks(input.as_str());
    println!(
        "Part 1. The message after executing move instructions is {}",
        execute_move_instructions(input.as_str(), &mut stacks.clone(), true)
    );
    println!(
        "Part 2. The message after executing move instructions is {}",
        execute_move_instructions(input.as_str(), &mut stacks.clone(), false)
    );
}

fn parse_stacks(input: &str) -> Vec<Vec<char>> {
    let mut stacks = vec![vec![]; 9];
    for line in input.lines() {
        if line.starts_with(" 1") {
            break;
        }
        for stack_index in 0..9 {
            let char_pos = stack_index * 4 + 1;
            let crate_marking = line.chars().nth(char_pos).unwrap();
            if crate_marking != ' ' {
                stacks[stack_index].push(crate_marking);
            }
        }
    }

    stacks.iter_mut().for_each(|stack| stack.reverse());

    stacks
}

fn execute_move_instructions(
    input: &str,
    stacks: &mut Vec<Vec<char>>,
    reverse_move: bool,
) -> String {
    input
        .lines()
        .skip_while(|line| !line.starts_with("move"))
        .for_each(|line| {
            let words: Vec<&str> = line.split(' ').collect();

            let amount: usize = words[1].parse().unwrap();
            let from: usize = words[3].parse().unwrap();
            let to: usize = words[5].parse().unwrap();

            let from_len = stacks[from - 1].len();
            let mut to_move: Vec<_> = if reverse_move {
                stacks[from - 1].drain(from_len - amount..).rev().collect()
            } else {
                stacks[from - 1].drain(from_len - amount..).collect()
            };
            stacks[to - 1].append(&mut to_move);
        });

    let mut result: String = "".to_string();
    stacks
        .iter()
        .for_each(|stack| result.push(*stack.last().unwrap()));
    result
}
