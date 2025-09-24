pub mod data;
pub mod formats;
/// File export module for saving results to files
///
/// This module handles all file export functionality, separating it from
/// the output module which handles console display operations.
///
/// - `trait_impl`: Export trait and format enum for unified export interface
/// - `formats`: Concrete implementations for CSV and JSON exporters
/// - `data`: Data structures specifically for serialization and export
pub mod trait_impl;

pub use data::ExportData;
pub use trait_impl::{export_data, ExportFormat};
