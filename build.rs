fn main() -> Result<(), Box<dyn std::error::Error>> {
    // tonic_build::compile_protos("apis/kennylong/v1/webhook.proto")?;
    tonic_build::configure()
        .file_descriptor_set_path("target/api_descriptor.pb")
        .compile(&["apis/github/v1/webhook.proto"], &["apis", "googleapis"])?;
    Ok(())
}
