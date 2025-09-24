pub fn parse_date_string(date_str: &str) -> Option<(String, String, String)> {
    let parts: Vec<&str> = date_str.split(['-', '/']).collect();
    if parts.len() == 3 {
        Some((
            parts[0].to_string(),
            parts[1].to_string(),
            parts[2].to_string(),
        ))
    } else {
        None
    }
}
