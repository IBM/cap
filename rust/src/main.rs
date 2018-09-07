extern crate reqwest;
extern crate atom_syndication;

use std::io::BufReader;

use atom_syndication::Feed;

fn main() {
    let resp = reqwest::get("https://alerts.weather.gov/cap/us.php?x=1").unwrap();

    let reader = BufReader::new(resp);

    let feed = Feed::read_from(BufReader::new(reader)).unwrap();

    println!("{}", feed.to_string());
}
