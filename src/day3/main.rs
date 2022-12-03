use std::fs;

fn main() {
    let input = fs::read_to_string("src/day3/input.txt").expect("Should be able to read the input");
    println!(
        "Task 1. The sum of the common item priorities is {}",
        sum_of_priorities(input.as_str()),
    );
    println!(
        "Task 2. The sum of all badge priorities is {}",
        sum_of_badge_priorities(input.as_str()),
    );
}

fn sum_of_priorities(input: &str) -> u64 {
    input.lines().fold(0, |acc, line| {
        let middle_index = line.len() / 2;
        let common_item = find_common_item(&line[..middle_index], &line[middle_index..]);
        acc + item_score(common_item)
    })
}

fn sum_of_badge_priorities(input: &str) -> u64 {
    let mut sum = 0;
    for chunk in input.lines().collect::<Vec<_>>().chunks_exact(3) {
        let badge = find_common_item_in_three_collections(chunk[0], chunk[1], chunk[2]);
        sum += item_score(badge);
    }
    sum
}

fn find_common_item(compartment1: &str, compartment2: &str) -> char {
    for c in compartment1.chars() {
        if compartment2.find(c).is_some() {
            return c;
        }
    }
    panic!("there should always be a common item in both compartments");
}

fn find_common_item_in_three_collections(
    collection1: &str,
    collection2: &str,
    collection3: &str,
) -> char {
    for c1 in collection1.chars() {
        for c2 in collection2.chars() {
            for c3 in collection3.chars() {
                if c1 == c2 && c2 == c3 {
                    return c1;
                }
            }
        }
    }
    panic!("there should always be a common item in the three compartments");
}

fn item_score(item: char) -> u64 {
    if item.is_uppercase() {
        return item as u64 - 65 + 27;
    }
    item as u64 - 97 + 1
}
