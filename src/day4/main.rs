use std::fs;

fn main() {
    let input = fs::read_to_string("src/day4/input.txt").expect("Should be able to read input");
    let range_pairs = parse_range_pairs(input.as_str());
    println!(
        "Part 1. There are {} assignments pairs where one range is fully covered by the other",
        fully_covered_ranges_count(&range_pairs)
    );
    println!(
        "Part 2. There are {} assignments pairs where one range has overlap with the other",
        ranges_with_overlap_count(&range_pairs)
    );
}

fn parse_range_pairs(input: &str) -> Vec<(Range, Range)> {
    input
        .lines()
        .map(|line| {
            let (assignment1, assignment2) = line.split_once(',').unwrap();
            let assignment1_ranges = assignment1.split_once('-').unwrap();
            let range1 = Range {
                lower: assignment1_ranges.0.parse().unwrap(),
                upper: assignment1_ranges.1.parse().unwrap(),
            };
            let assignment2_ranges = assignment2.split_once('-').unwrap();
            let range2 = Range {
                lower: assignment2_ranges.0.parse().unwrap(),
                upper: assignment2_ranges.1.parse().unwrap(),
            };
            (range1, range2)
        })
        .collect()
}

fn fully_covered_ranges_count(range_pairs: &Vec<(Range, Range)>) -> usize {
    range_pairs
        .iter()
        .filter(|(range1, range2)| {
            range1.is_fully_covered_by(&range2) || range2.is_fully_covered_by(&range1)
        })
        .count()
}

fn ranges_with_overlap_count(range_pairs: &Vec<(Range, Range)>) -> usize {
    range_pairs
        .iter()
        .filter(|(range1, range2)| range1.has_overlap(&range2))
        .count()
}

struct Range {
    lower: u8,
    upper: u8,
}

impl Range {
    fn is_fully_covered_by(&self, other: &Range) -> bool {
        self.lower >= other.lower && self.upper <= other.upper
    }

    fn has_overlap(&self, other: &Range) -> bool {
        self.lower >= other.lower && self.lower <= other.upper || // Lower is inside the other's range.
        self.upper >= other.lower && self.upper <= other.upper || // Upper is inside the other's range.
        self.lower < other.lower && self.upper > other.upper  || // Self fully contains the other's range.
        other.lower < self.lower && other.upper > self.upper // Other fully contains the self's range.
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn fully_covered_ranges_count_against_examples() {
        let example = "2-4,6-8\n2-3,4-5\n5-7,7-9\n2-8,3-7\n6-6,4-6\n2-6,4-8";
        let res = fully_covered_ranges_count(&parse_range_pairs(example));
        println!("{}", res);
        assert!(res == 2);
    }
}
