use tonic::{Request, Response, Status};
use v1::webhook_server::Webhook;
use v1::CreateWebhookRequest;

pub mod v1 {
    tonic::include_proto!("github.v1"); // The string specified here must match the proto package name
}

#[derive(Debug, Default)]
pub struct Service {
    secret: String,
}

#[tonic::async_trait]
impl Webhook for Service {
    async fn create_webhook(
        &self,
        request: Request<CreateWebhookRequest>,
    ) -> Result<Response<()>, Status> {
        println!("headers: {:?}", request.metadata());
        if let Some(ev) = request.metadata().get("x-gitHub-event") {
            let ev = ev.to_str().unwrap();
            if ev != "push" {
                println!("unhandled event {}", ev);
                return Err(Status::invalid_argument(format!(
                    "expect event to be push, got {ev}"
                )));
            }
        } else {
            return Err(Status::invalid_argument("no X-GitHub-Event header found"));
        }
        println!("payload: {:#?}", request.into_inner());
        // TODO: ghost admin api...
        Ok(Response::new(()))
    }
}
