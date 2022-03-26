use chrono::Duration;
use serde::{Deserialize, Serialize};
use std::{error, ops::Add};

/// Jwt fetches the content api token.
struct Jwt {
    key_id: Option<String>,
    secret: String,
    claims: Claims,
}

#[derive(Debug, Serialize, Deserialize)]
struct Claims {
    aud: String, // Optional. Audience
    exp: usize, // Required (validate_exp defaults to true in validation). Expiration time (as UTC timestamp)
    iat: usize, // Optional. Issued at (as UTC timestamp)
                // iss: String, // Optional. Issuer
                // nbf: usize, // Optional. Not Before (as UTC timestamp)
                // sub: String, // Optional. Subject (whom token refers to)
}

impl Jwt {
    fn from_ghost_key(key: &str) -> Self {
        let now = chrono::Utc::now();
        let strs: Vec<&str> = key.split(":").collect();
        Jwt {
            key_id: Some(strs[0].to_owned()),
            secret: strs[1].to_owned(),
            claims: Claims {
                aud: "/v3/admin/".to_owned(),
                exp: now.add(Duration::minutes(5)).timestamp() as usize, // ghost api only have 5min...
                iat: now.timestamp() as usize,
            },
        }
    }

    async fn get_token(&self) -> Result<String, Box<dyn error::Error>> {
        let mut hdr = jsonwebtoken::Header::default();
        hdr.kid = self.key_id.clone();

        let secret = hex::decode(&self.secret)?;
        let token = jsonwebtoken::encode(
            &hdr,
            &self.claims,
            &jsonwebtoken::EncodingKey::from_secret(&secret),
        )?;
        Ok(token)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test() {
        // this is just a local test id.
        let j = Jwt::from_ghost_key("623ebe1d55fb66015e7e6902:5e7d7a833f2063d47953c27627703a7379f2309230b6f47674508799cd120465");
        let res = j.get_token().await;
        assert_eq!(res.is_ok(), true);
        println!("{}", res.unwrap());
    }
}
