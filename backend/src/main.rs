use actix_cors::Cors;
use actix_web::{web, App, HttpServer};
use sqlx::postgres::PgPoolOptions;

mod config;
mod handlers;
mod models;

use crate::config::Settings;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let settings = Settings::new().expect("Failed to load settings");

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&settings.database_url)
        .await
        .expect("Failed to create pool");

    HttpServer::new(move || {
        let cors = Cors::default()
            .allow_any_origin()
            .allow_any_method()
            .allow_any_header()
            .max_age(3600);

        App::new()
            .wrap(cors)
            .app_data(web::Data::new(pool.clone()))
            .route("/shorten", web::post().to(handlers::create_short_url))
            .route("/{short_code}", web::get().to(handlers::redirect))
    })
    .bind(settings.server_addr)?
    .run()
    .await
}
