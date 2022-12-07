use std::fs;

// I got too annoyed and switched to Go this day. I'll commit this anyway...

fn main() {
    let input = fs::read_to_string("src/day7/input.txt").expect("Should be able to read the input");
    parse_directories(&input);
}

fn parse_directories(input: &str) -> Directory {
    let mut root = Directory::new(String::from("/"));
    // let mut current_directory = &mut root;
    let mut directories: Vec<&mut Directory> = vec![&mut root];
    for line in input.lines() {
        // if line.starts_with("$ ls") {
        // continue;
        if line.starts_with("$ cd") {
            let target = line.strip_prefix("$ cd ").unwrap();
            if target == ".." {
                directories.pop();
            } else {
                // directories.push(current_directory);
                directories.push(
                    directories
                        .last_mut()
                        .unwrap()
                        .subdirectories
                        .iter_mut()
                        .find(|d| d.name == target)
                        .unwrap(),
                )
            }
        } else if !line.starts_with('$') {
            if line.starts_with("dir") {
                directories
                    .last_mut()
                    .unwrap()
                    .subdirectories
                    .push(Directory::new(line.split_once(' ').unwrap().1.to_string()));
            } else {
                let (size, name) = line.split_once(' ').unwrap();
                directories
                    .last_mut()
                    .unwrap()
                    .files
                    .push(File::new(name.to_string(), size.parse().unwrap()))
            }
        }
    }

    root
}

struct Directory {
    name: String,
    subdirectories: Vec<Directory>,
    files: Vec<File>,
}

struct File {
    name: String,
    size: u64,
}

impl Directory {
    fn new(name: String) -> Directory {
        Directory {
            name,
            subdirectories: vec![],
            files: vec![],
        }
    }

    fn size(&self) -> u64 {
        self.files.iter().map(|file| file.size).sum()
    }
}

impl File {
    fn new(name: String, size: u64) -> File {
        File { name, size }
    }
}
