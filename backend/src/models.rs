use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct Url {
    pub id: i32,
    pub short_code: String,
    pub long_url: String,
}

#[derive(Deserialize)]
pub struct CreateUrl {
    pub long_url: String,
}
