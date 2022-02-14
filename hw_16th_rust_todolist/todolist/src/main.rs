use actix_files::NamedFile;
use actix_web::{
    delete, get, patch, post, web, App, HttpRequest, HttpResponse, HttpServer, Responder, Result,
};
use serde::Deserialize;
use std::path::PathBuf;

#[derive(Deserialize)]
struct SubjectRequest {
    page: String,
    group: String,
}

#[derive(Deserialize)]
struct NewSubjectRequest {
    subject: String,
}

#[derive(Deserialize)]
struct UpdateSubjectRequest {
    id: String,
    subject: String,
    status: String,
}

#[derive(Debug, Deserialize)]
pub enum ResponseType {
    id,
    status,
    subject,
}

async fn index(_req: HttpRequest) -> Result<NamedFile> {
    let path: PathBuf = "./index.html".parse().unwrap();
    Ok(NamedFile::open(path)?)
}

#[get("/api")]
async fn get(web::Query(request): web::Query<SubjectRequest>) -> impl Responder {
    println!("{},{}", request.page, request.group);
    HttpResponse::Ok().body("Hello world!")
}

#[post("/api")]
async fn post(new_subject_request: web::Json<NewSubjectRequest>) -> impl Responder {
    HttpResponse::Ok().body("")
}

#[patch("/api")]
async fn patch(req_body: String) -> impl Responder {
    HttpResponse::Ok().body(req_body)
}

#[delete("/api")]
async fn delete(req_body: String) -> impl Responder {
    HttpResponse::Ok().body(req_body)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(index))
            .service(get)
            .service(post)
            .service(patch)
            .service(delete)
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
