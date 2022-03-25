# kennylong.io

The source code of kennylong.io.

Working in progress, currently rewrite with Rust.

> For the v3 version which written in Go, see branch master.

## Features

Using the Github webhook(or Github app?) to receive event(push, pr, etc.), then rendering the blog to it's backend.

The blog engine backend should be abstracted, currently I am using ghost.

Finally, it's cloud native.