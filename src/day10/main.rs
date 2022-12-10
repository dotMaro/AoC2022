use std::fs;

// I got too annoyed again and remade it in Go.

fn main() {
    let input = fs::read_to_string("src/day10/input.txt").unwrap();
}

struct CPU {
    x: i64,
    cycle: usize,
    current_instruction: Option<Instruction>,
    instructions: Vec<Instruction>,
    pc: usize,
}

struct Instruction {
    instruction_type: InstructionType,
    value: i8,
    cycles_remaining: u8,
}

impl Instruction {
    fn new(instruction_type: InstructionType, value: Option<i8>) -> Self {
        let instruction_value = match value {
            Some(v) => v,
            None => 0,
        };
        let cycles_remaining = instruction_type.cycles_required();
        Instruction {
            instruction_type,
            value: instruction_value,
            cycles_remaining,
        }
    }
}

enum InstructionType {
    Noop,
    Addx,
}

impl InstructionType {
    fn new(s: &str) -> Self {
        match s {
            "noop" => InstructionType::Noop,
            "addx" => InstructionType::Addx,
            _ => panic!(),
        }
    }

    fn cycles_required(&self) -> u8 {
        match self {
            InstructionType::Noop => 1,
            InstructionType::Addx => 2,
        }
    }
}

impl CPU {
    fn new(input: &str) -> Self {
        let instructions = input
            .lines()
            .map(|l| {
                let words: Vec<&str> = l.split(' ').collect();
                let value = if words.len() > 1 {
                    Some(words[1].parse().unwrap())
                } else {
                    None
                };
                Instruction::new(InstructionType::new(words[0]), value)
            })
            .collect();
        CPU {
            x: 1,
            cycle: 0,
            current_instruction: None,
            instructions,
            pc: 0,
        }
    }

    fn cycle(&mut self, input: &str) {
        let cpu = self;
        if let Some(instruction) = cpu.current_instruction {
            cpu.current_instruction.unwrap().cycles_remaining -= 1;
            if instruction.cycles_remaining == 0 {
                cpu.execute_instruction(&instruction);
            }
        }
    }

    fn execute_instruction(&mut self, instruction: &Instruction) {
        match instruction.instruction_type {
            InstructionType::Addx => {
                self.x += instruction.value as i64;
            }
            InstructionType::Noop => {}
        }
    }
}
