use anyhow::Result;
use clap::{Parser, Subcommand};
use console::{style, Emoji};

mod commands;
mod export;
mod output;
mod progress;
mod utils;

use commands::{handle_duplicates, handle_organize, DuplicatesArgs, OrganizeArgs};

static LOOKING_GLASS: Emoji = Emoji("🔍 ", "");
static FILES: Emoji = Emoji("📁 ", "");
static DUPLICATE: Emoji = Emoji("🔄 ", "");
static WARNING: Emoji = Emoji("⚠️ ", "");

#[derive(Parser)]
#[command(name = "image-manager-cli")]
#[command(about = "A CLI tool for image organization and duplicate detection")]
#[command(version = "0.1.0")]
#[command(author = "Elias Schneider")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Preview how images would be organized by date
    Organize(OrganizeArgs),
    /// Find duplicate images in a directory
    Duplicates(DuplicatesArgs),
}

fn main() {
    let cli = Cli::parse();

    match run(cli) {
        Ok(_) => {
            println!("\n{}", style("✓ Operation completed successfully").green());
        }
        Err(e) => {
            eprintln!("\n{} {}", WARNING, style(format!("Error: {}", e)).red());
            std::process::exit(1);
        }
    }
}

fn run(cli: Cli) -> Result<()> {
    match cli.command {
        Commands::Organize(args) => {
            println!(
                "{} {} Scanning directory for organization preview...",
                LOOKING_GLASS,
                style("Organize").cyan()
            );
            handle_organize(args)
        }
        Commands::Duplicates(args) => {
            println!(
                "{} {} Scanning directory for duplicates...",
                LOOKING_GLASS,
                style("Duplicates").cyan()
            );
            handle_duplicates(args)
        }
    }
}
