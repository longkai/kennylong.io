mod apis;
mod jwt;

use apis::github;
use tonic::transport::Server;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::]:50051".parse()?;
    let github = github::Service::default();

    Server::builder()
        .add_service(github::v1::webhook_server::WebhookServer::new(github))
        .serve(addr)
        .await?;

    Ok(())
}
