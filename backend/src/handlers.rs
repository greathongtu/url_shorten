use crate::models::{CreateUrl, Url};
use actix_web::{web, HttpResponse, Responder};
use rand::Rng;
use sqlx::PgPool;

pub async fn create_short_url(
    pool: web::Data<PgPool>,
    url: web::Json<CreateUrl>,
) -> impl Responder {
    let short_code = generate_short_code();

    let result = sqlx::query_as!(
        Url,
        "INSERT INTO urls (short_code, long_url) VALUES ($1, $2) RETURNING *",
        short_code,
        url.long_url
    )
    .fetch_one(pool.get_ref())
    .await;

    match result {
        Ok(url) => HttpResponse::Ok().json(url),
        Err(_) => HttpResponse::InternalServerError().finish(),
    }
}

pub async fn redirect(pool: web::Data<PgPool>, short_code: web::Path<String>) -> impl Responder {
    let result = sqlx::query_as!(
        Url,
        "SELECT * FROM urls WHERE short_code = $1",
        short_code.into_inner()
    )
    .fetch_optional(pool.get_ref())
    .await;

    match result {
        Ok(Some(url)) => HttpResponse::Found()
            .append_header(("Location", url.long_url))
            .finish(),
        Ok(None) => HttpResponse::NotFound().finish(),
        Err(_) => HttpResponse::InternalServerError().finish(),
    }
}

fn generate_short_code() -> String {
    const CHARSET: &[u8] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    const SHORT_CODE_LEN: usize = 6;
    let mut rng = rand::thread_rng();

    (0..SHORT_CODE_LEN)
        .map(|_| {
            let idx = rng.gen_range(0..CHARSET.len());
            CHARSET[idx] as char
        })
        .collect()
}
