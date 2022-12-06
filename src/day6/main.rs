use std::collections::HashSet;
use std::fs;

fn main() {
    let input = fs::read_to_string("src/day6/input.txt").expect("Should be able to read the input");

    println!(
        "Task 1. The marker after the first four unique characters is at {}",
        index_of_first_sequence_with_unique_characters(&input, 4)
    );
    println!(
        "Task 2. The marker after the first 14 unique characters is at {}",
        index_of_first_sequence_with_unique_characters(&input, 14)
    );
}

fn index_of_first_sequence_with_unique_characters(input: &str, sequence_len: usize) -> usize {
    input
        .chars()
        .collect::<Vec<char>>()
        .windows(sequence_len)
        .enumerate()
        .find(|(_, w)| !has_duplicates(w))
        .map(|(index, _)| index + sequence_len)
        .unwrap()
}

fn has_duplicates(chars: &[char]) -> bool {
    let mut previous_chars: HashSet<char> = HashSet::new();
    for c in chars {
        if previous_chars.contains(c) {
            return true;
        }
        previous_chars.insert(*c);
    }
    false
}

#[cfg(test)]
mod tests {
    use crate::index_of_first_sequence_with_unique_characters;

    #[test]
    fn examples_part1() {
        assert!(
            index_of_first_sequence_with_unique_characters("bvwbjplbgvbhsrlpgdmjqwftvncz", 4) == 5
        )
    }

    #[test]
    fn examples_part2() {
        assert!(
            index_of_first_sequence_with_unique_characters("mjqjpqmgbljsphdztnvjfqwrcgsmlb", 14)
                == 19
        )
    }
}
