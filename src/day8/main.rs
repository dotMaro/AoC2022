use std::fs;

fn main() {
    let input = fs::read_to_string("src/day8/input.txt").expect("Should be able to read the input");
    let trees = parse_trees(&input);
    println!("Part 1. {}", visible_tree_count(&trees));
    println!("Part 2. {}", highest_scenic_score(&trees));
}

fn parse_trees(input: &str) -> Vec<Vec<u8>> {
    input
        .lines()
        .map(|line| {
            line.chars()
                .map(|c| c.to_digit(10).unwrap() as u8)
                .collect()
        })
        .collect()
}

fn visible_tree_count(trees: &Vec<Vec<u8>>) -> usize {
    let mut count = 0;
    let width = trees[0].len();
    for y in 0..trees.len() {
        for x in 0..width {
            if tree_is_visible(trees, x, y) {
                count += 1;
            }
        }
    }
    count
}

fn highest_scenic_score(trees: &Vec<Vec<u8>>) -> usize {
    let mut highest = 0;
    let width = trees[0].len();
    for y in 0..trees.len() {
        for x in 0..width {
            let scenic_score = scenic_score(trees, x, y);
            if scenic_score > highest {
                highest = scenic_score;
            }
        }
    }
    highest
}

fn tree_is_visible(trees: &Vec<Vec<u8>>, x: usize, y: usize) -> bool {
    if x == 0 || x == trees[y].len() - 1 || y == 0 || y == trees.len() - 1 {
        return true;
    }

    let tree_height = trees[y][x];
    trees[y][..x].iter().all(|t| *t < tree_height)
        || trees[y][x + 1..].iter().all(|t| *t < tree_height)
        || trees[..y].iter().all(|row| row[x] < tree_height)
        || trees[y + 1..].iter().all(|row| row[x] < tree_height)
}

fn scenic_score(trees: &Vec<Vec<u8>>, x: usize, y: usize) -> usize {
    if x == 0 || x == trees[y].len() - 1 || y == 0 || y == trees.len() - 1 {
        return 0;
    }

    let tree_height = trees[y][x];
    let left_visible = trees[y][..x]
        .iter()
        .rev()
        .take_while(|t| **t < tree_height)
        .count();
    let left_score = if left_visible == x {
        left_visible
    } else {
        left_visible + 1
    };
    let right_visible = trees[y][x + 1..]
        .iter()
        .take_while(|t| **t < tree_height)
        .count();
    let right_score = if right_visible == trees[y].len() - x - 1 {
        right_visible
    } else {
        right_visible + 1
    };
    let up_visible = trees[..y]
        .iter()
        .rev()
        .take_while(|row| row[x] < tree_height)
        .count();
    let up_score = if up_visible == y {
        up_visible
    } else {
        up_visible + 1
    };
    let down_visible = trees[y + 1..]
        .iter()
        .take_while(|row| row[x] < tree_height)
        .count();
    let down_score = if down_visible == trees.len() - y - 1 {
        down_visible
    } else {
        down_visible + 1
    };
    left_score * right_score * up_score * down_score
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn example_part1() {
        let example = "30373
25512
65332
33549
35390";
        let trees = parse_trees(example);
        assert!(visible_tree_count(&trees) == 21);
    }

    #[test]
    fn example_part2() {
        let example = "30373
25512
65332
33549
35390";
        let trees = parse_trees(example);
        let score = highest_scenic_score(&trees);
        assert!(score == 8);
    }
}
